#!/bin/bash

TARGET_DIR="${1:-.}"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

shopt -s nullglob
MJML_FILES=("$TARGET_DIR"/*.mjml)

if [ ${#MJML_FILES[@]} -eq 0 ]; then
    echo -e "${RED}[ERROR]:${NC} There're no .mjml files in '$TARGET_DIR'"
    exit 1
fi

normalize_html_ids() {
    local file="$1"

    if [[ -z "$file" ]]; then
        echo -e "${RED}[ERROR]${NC} Usage: normalize_html_ids <file.html>" >&2
        return 1
    fi

    if [[ ! -f "$file" ]]; then
        echo -e "${RED}[ERROR]${NC} File '$file' not found." >&2
        return 1
    fi

    # extract generated ids
    local ids
    ids=$(
        grep -Eo 'id="[a-fA-F0-9]{16}"|mj-carousel-[a-fA-F0-9]{16}-icons-cell' "$file" | 
        grep -Eo '[a-fA-F0-9]{16}' | 
        awk '!seen[$0]++'
    )

    if [[ -z "$ids" ]]; then
        return 0
    fi

    local sed_cmds=""
    local counter=1
    local new_id

    for old_id in $ids; do
        new_id=$(printf "%016d" "$counter")
        sed_cmds+="s/$old_id/$new_id/g;"
        ((counter++))
    done

    local tmp_file="${file}.tmp"
    
    if sed -e "$sed_cmds" "$file" > "$tmp_file"; then
        mv "$tmp_file" "$file"
    else
        echo -e "${RED}[ERROR]${NC} Unable to process the file '$file'." >&2
        rm -f "$tmp_file"
        return 1
    fi
}

for file in "${MJML_FILES[@]}"; do
    filename=$(basename -- "$file")
    basename="${filename%.*}"

    echo -e "Testing: ${YELLOW}$filename${NC} ..."

    mjml_html="./mjml-${basename}.html"
    mgml_html="./mgml-${basename}.html"
    mjml_min="./mjml-${basename}.min.html"
    mgml_min="./mgml-${basename}.min.html"

    # generate via mjml
    mjml -r "$file" --config.beautify=false -o "$mjml_html" > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "${RED}[ERROR]${NC} The mjml utility exited with an error for $filename"
        continue
    fi

    # generate via mgml
    mgml "$file" "$mgml_html" > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "${RED}[ERROR]${NC} The mgml utility exited with an error for $filename"
        continue
    fi

    minify --type=html --html-keep-conditional-comments --html-keep-quotes --html-keep-end-tags --html-keep-document-tags -o "$mjml_min" "$mjml_html" > /dev/null 2>&1
    minify --type=html --html-keep-conditional-comments --html-keep-quotes --html-keep-end-tags --html-keep-document-tags -o "$mgml_min" "$mgml_html" > /dev/null 2>&1

    normalize_html_ids "$mjml_min"
    normalize_html_ids "$mgml_min"

    diff_output=$(diff -u "$mjml_min" "$mgml_min")
    diff_status=$?

    if [ $diff_status -eq 0 ]; then
        echo -e "${GREEN}[SUCCESS]${NC} There're no differences"
        rm -f "$mjml_html" "$mgml_html" "$mjml_min" "$mgml_min"
    else
        echo -e "${RED}[ERROR]${NC} Differences were found in '$filename'"
    fi
done
