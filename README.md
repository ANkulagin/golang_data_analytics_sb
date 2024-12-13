![Go](https://img.shields.io/badge/Go-1.23-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Version](https://img.shields.io/badge/Version-v1.0.0-blue.svg)

# Golang Markdown Converter Second Brain


## Other Languages

This documentation is also available in:
- [Russian](README.ru.md)

## Description

**Golang Markdown Converter Second Brain** is a Go-based application designed to automate the conversion of Obsidian notes into HTML format. It ensures efficient transformation while preserving directory structures and metadata, making note management more convenient and organized.

## Key Features

- **Automatic Conversion**: Converts all Markdown files in a specified directory into HTML.
- **Directory Structure Preservation**: Maintains the folder hierarchy of the source files.
- **Metadata Integration**: Embeds YAML Front Matter metadata as HTML comments (not individual attributes).
- **Flexible Logging Levels**: Uses the `logrus` library for adjustable log verbosity.
- **Change Detection**: Avoids unnecessary HTML file overwriting by checking for file modifications.

## Project Structure and Visual Representation
- [Flowchart](docs/Flowchart.mmd)
- [Class Diagram](docs/ClassDiagram.mmd)

## Installation

### Prerequisites

- **Go**: Ensure Go version 1.23 or higher is installed. Download and install Go from the [official website](https://golang.org/dl/).

### Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/ANkulagin/golang_markdown_converter_sb.git
cd golang_markdown_converter_sb
```

### Install Dependencies

Ensure all dependencies are installed:

```bash
go mod tidy
```

### Configuration

The application configuration is stored in the `configs/config.yaml` file. You can adjust the following settings:

```yaml
src_dir: "/home/ankul/obsidian/_notes/daily"   
dest_dir: "/home/ankul/_html/daily"          
log_level: "info"                         
```

#### Configuration Parameters

- `src_dir`: The absolute or relative path to the directory containing your Markdown files. The application will recursively find all `.md` files in this directory.
- `dest_dir`: The absolute or relative path to the directory where generated HTML files will be saved. Directory structure will be preserved.
- `log_level`: The logging level of the application. Possible values:
  - trace
  - debug
  - info
  - warn
  - error
  - fatal
  - panic

## Usage

### Running the Application

To run the converter, use the following command (the flag is optional if the default configuration file is used):

```bash
go run cmd/daily/main.go -config=configs/config.yaml
```

### Building an Executable

You can also build the application into an executable file:

```bash
go build -o markdown_converter cmd/daily/main.go
```

Afterward, you can run the application without recompiling:

```bash
./markdown_converter -config=configs/config.yaml
```

### Testing

To run the tests, use the following command:

```bash
go test ./internal/service/converter/...
```

## Contributing and Support

Contributions to the project are welcome! If you want to make changes or add new features, feel free to fork the repository and submit your changes on [GitHub](https://github.com/ANkulagin/golang_markdown_converter_sb).

## Contact

If you have any questions or suggestions, feel free to reach out:
- Telegram: [ANkulagin03](https://t.me/ANkulagin03)

## License

This project is licensed under the MIT License. For more details, see the [LICENSE](LICENSE) file.

---


