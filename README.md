# Uploading excel spreedsheet using Go
Framework: gin

## cURL example

```bash
curl --location '0.0.0.0:8080/upload' \
--form 'file=@"{...}/excel_sheet.xlsx"'
```