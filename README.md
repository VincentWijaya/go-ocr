# go-ocr
This is my golang service template that implement clean architecture.

Feel free to use it 😊

---
# How to run:

```
1. Clone this repo for sure 😊
2. Rename config.yaml.example to ocr.yaml in ./files/etc/ocr/
3. go mod tidy
4. go run cmd/main.go
```

# Project structer:
```
📦go-ocr
 ┣ 📂cmd
 ┃ ┗ 📜main.go
 ┣ 📂constant
 ┃ ┗ 📂errs
 ┃ ┃ ┗ 📜errors.go
 ┣ 📂files
 ┃ ┗ 📂etc
 ┃ ┃ ┗ 📂ocr
 ┃ ┃ ┃ ┣ 📜config.yaml.example
 ┃ ┃ ┃ ┗ 📜ocr.yaml
 ┣ 📂internal
 ┃ ┣ 📂app
 ┃ ┃ ┣ 📂handler
 ┃ ┃ ┃ ┣ 📜helper.go
 ┃ ┃ ┃ ┗ 📜init.go
 ┃ ┃ ┣ 📂repo
 ┃ ┃ ┗ 📂usecase
 ┃ ┗ 📂entity
 ┃ ┃ ┗ 📜entity.go
 ┣ 📂log
 ┃ ┗ 📜ocr.info.log
 ┣ 📜.gitignore
 ┣ 📜README.md
 ┣ 📜go.mod
 ┗ 📜go.sum
```

---
