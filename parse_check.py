with open("core/search_parser_test.go", "r") as f:
    lines = f.readlines()
data_lines = []
in_data = False
for line in lines:
    if line.startswith("var sampleData = `"):
        in_data = True
    elif line.startswith("`") and in_data:
        in_data = False
    elif in_data:
        if line.startswith("!"):
            data_lines.append(line.strip())

print(f"Total lines: {len(data_lines)}")
for i, l in enumerate(data_lines):
    print(f"{i+1}: {l}")
