import re

regex = re.compile(r"^!(?P<server>[^\s]+)\s+(?:%[A-Za-z0-9]+%\s+)?(?:(?P<author>[^-]+?)\s+-\s+)?(?P<title>.*?\.[a-zA-Z0-9]+)\.(?P<format>[a-zA-Z0-9]+)(?:\s+::INFO::\s*(?P<size>[^\s]+))?")

def parseLineRegex(line):
    m = regex.match(line)
    if m:
        return m.groupdict()
    return None

lines = [
	"!dragnbreaker Fitzgerald, F Scott - Novel 03 - The Great Gatsby (retail).epub  ::INFO:: 1.7MB",
	"!peapod The Great Gatsby.pdf  ::INFO:: 254.73KB",
	"!DeathCookie Emma_L_Adams_Heritage_of_Fire_03_Inferno.epub.rar  ::INFO:: 530.0KB",
	"!Horla F Scott Fitzgerald - The Great Gatsby (retail) (epub).epub",
	"!FWServer %F77FE9FF1CCD% Michael Haag - Inferno Decoded.epub  ::INFO:: 8.00MB",
    "!Horla Linda Howard -[Raintree 01]- Inferno.doc",
]

for l in lines:
    print(l, "->", parseLineRegex(l))
