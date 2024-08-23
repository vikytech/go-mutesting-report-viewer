# Go Mutesting HTML Report Viewer

![screenshot](./assets/img/screenshot.png)

## To directly run the binary
```zsh
./go_mutesting_html_report -file <PATH_TO_JSON_REPORT>
```
> Note: Make sure you also have the template.html in the same directory as the executable

If you have trouble running on Mac os run: `xattr -d com.apple.quarantine go_mutesting_html_report`

## To run locally
```zsh
go run go_mutesting_html_report.go -file <PATH_TO_JSON_REPORT>
```

### Features
- View Diff with changes highlighted line-by-line
- Mark as Viewed and collapse the same

### Upcoming features
 - Filter
 - Enhanced UI / UX
 - Compare Reports

### Dependencies
- [diff2html](https://diff2html.xyz)