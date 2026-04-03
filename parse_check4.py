def check_title_author(line, fileTypes):
    server = ""
    author = ""
    title = ""
    format = ""
    size = "N/A"

    firstSpace = line.find(" ")
    server = line[1:firstSpace]

    rest = line[firstSpace+1:]
    infoIdx = rest.find(" ::INFO:: ")
    if infoIdx != -1:
        size = rest[infoIdx + len(" ::INFO:: "):].split()[0]
        rest = rest[:infoIdx]

    lineLower = rest.lower()

    for ext in fileTypes:
        endTitle = lineLower.rfind("." + ext)
        if endTitle == -1:
            continue

        format = ext
        if ext in ("rar", "zip"):
            for ext2 in fileTypes[:-2]:
                if ext2 in lineLower[:endTitle]:
                    format = ext2

        dashChar = rest.find(" - ")
        if dashChar == -1:
            author = ""
            title = rest[:endTitle]
        else:
            author = rest[:dashChar]
            if "%" in author:
                split = author.split(" ", 1)
                if len(split) == 2:
                    author = split[1]
            title = rest[dashChar + len(" - ") : endTitle]
        break

    return f"server: {server}, author: {author}, title: {title}, format: {format}, size: {size}"

fileTypes = [
	"epub",
	"mobi",
	"azw3",
	"html",
	"rtf",
	"pdf",
	"cdr",
	"lit",
	"cbr",
	"doc",
	"htm",
	"jpg",
	"txt",
	"rar", # Compressed extensions should always be last 2 items
	"zip",
]

lines = [
	"!dragnbreaker Fitzgerald, F Scott - Novel 03 - The Great Gatsby (retail).epub  ::INFO:: 1.7MB",
	"!peapod The Great Gatsby.pdf  ::INFO:: 254.73KB",
	"!DeathCookie Emma_L_Adams_Heritage_of_Fire_03_Inferno.epub.rar  ::INFO:: 530.0KB",
	"!Horla F Scott Fitzgerald - The Great Gatsby (retail) (epub).epub",
	"!FWServer %F77FE9FF1CCD% Michael Haag - Inferno Decoded.epub  ::INFO:: 8.00MB",
    "!Horla Linda Howard -[Raintree 01]- Inferno.doc",
]

for l in lines:
    print(check_title_author(l, fileTypes))
