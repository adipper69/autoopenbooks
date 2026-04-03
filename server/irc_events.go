package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
)

func (server *server) NewIrcEventHandler(client *Client) core.EventHandler {
	handler := core.EventHandler{}
	handler[core.SearchResult] = client.searchResultHandler(server.config.DownloadDir)
	handler[core.BookResult] = client.bookResultHandler(server.config.DownloadDir, server.config.DisableBrowserDownloads)
	handler[core.NoResults] = client.noResultsHandler
	handler[core.BadServer] = client.badServerHandler
	handler[core.SearchAccepted] = client.searchAcceptedHandler
	handler[core.MatchesFound] = client.matchesFoundHandler
	handler[core.Ping] = client.pingHandler
	handler[core.ServerList] = client.userListHandler(server.repository)
	handler[core.Version] = client.versionHandler(server.config.UserAgent)
	return handler
}

// searchResultHandler downloads from DCC server, parses data, and sends data to client
func (c *Client) searchResultHandler(downloadDir string) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(filepath.Join(downloadDir, "books"), text, nil)
		if err != nil {
			c.log.Println(err)
			select {
			case c.send <- newErrorResponse("Error when downloading search results."):
			case <-c.ctx.Done():
				c.log.Println("Client disconnected, dropping error response.")
			}
			return
		}

		bookResults, parseErrors, err := core.ParseSearchFile(extractedPath)
		if err != nil {
			c.log.Println(err)
			select {
			case c.send <- newErrorResponse("Error when parsing search results."):
			case <-c.ctx.Done():
				c.log.Println("Client disconnected, dropping error response.")
			}
			return
		}

		if len(bookResults) == 0 && len(parseErrors) == 0 {
			c.noResultsHandler(text)
			return
		}

		// Output all errors so parser can be improved over time
		if len(parseErrors) > 0 {
			c.log.Printf("%d Search Result Parsing Errors\n", len(parseErrors))
			for _, err := range parseErrors {
				c.log.Println(err)
			}
		}

		c.log.Printf("Sending %d search results.\n", len(bookResults))
		select {
		case c.send <- newSearchResponse(bookResults, parseErrors):
		case <-c.ctx.Done():
			c.log.Println("Client disconnected, dropping search response.")
		}

		err = os.Remove(extractedPath)
		if err != nil {
			c.log.Printf("Error deleting search results file: %v", err)
		}
	}
}

// bookResultHandler downloads the book file and sends it over the websocket
func (c *Client) bookResultHandler(downloadDir string, disableBrowserDownloads bool) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(filepath.Join(downloadDir, "books"), text, nil)
		if err != nil {
			c.log.Println(err)
			select {
			case c.send <- newErrorResponse("Error when downloading book."):
			case <-c.ctx.Done():
				c.log.Println("Client disconnected, dropping error response.")
			}
			return
		}

		c.log.Printf("Sending book entitled '%s'.\n", filepath.Base(extractedPath))
		select {
		case c.send <- newDownloadResponse(extractedPath, disableBrowserDownloads):
		case <-c.ctx.Done():
			c.log.Println("Client disconnected, dropping download response.")
		}
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (c *Client) noResultsHandler(_ string) {
	select {
	case c.send <- newErrorResponse("No results found for the query."):
	case <-c.ctx.Done():
		c.log.Println("Client disconnected, dropping error response.")
	}
}

// BadServer is called when the requested download fails because the server is not available
func (c *Client) badServerHandler(_ string) {
	select {
	case c.send <- newErrorResponse("Server is not available. Try another one."):
	case <-c.ctx.Done():
		c.log.Println("Client disconnected, dropping error response.")
	}
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (c *Client) searchAcceptedHandler(_ string) {
	select {
	case c.send <- newStatusResponse(NOTIFY, "Search accepted into the queue."):
	case <-c.ctx.Done():
		c.log.Println("Client disconnected, dropping status response.")
	}
}

// MatchesFound is called when the server finds matches for the user's query
func (c *Client) matchesFoundHandler(num string) {
	select {
	case c.send <- newStatusResponse(NOTIFY, fmt.Sprintf("Found %s results for your query.", num)):
	case <-c.ctx.Done():
		c.log.Println("Client disconnected, dropping status response.")
	}
}

func (c *Client) pingHandler(serverUrl string) {
	c.irc.Pong(serverUrl)
}

func (c *Client) versionHandler(version string) core.HandlerFunc {
	return func(line string) {
		c.log.Printf("Sending CTCP version response: %s", line)
		core.SendVersionInfo(c.irc, line, version)
	}
}

func (c *Client) userListHandler(repo *Repository) core.HandlerFunc {
	return func(text string) {
		repo.servers = core.ParseServers(text)
	}
}
