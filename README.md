# LINe ENumerator 

This Go program is a simple tool designed to recursively scan directories, identify source code files, and count the number of comment lines, empty lines, and code lines in each file type. It supports multiple file extensions and handles both single-line and multi-line comments.

## Features

- **Recursive Directory Scanning:** The program will search through all subdirectories starting from the provided root directory.
- **Source Code Analysis:** It differentiates between comments, empty lines, and code lines.
- **Multi-Language Support:** By recognizing different file extensions, it can handle various programming languages.

## Usage

### Prerequisites

- Go 1.16 or later installed.

### Installation

1. Clone the repository or download the source code.
2. Navigate to the directory where the source code is located.

### Running the Program

You can run the program using the following command:

```bash
go run linen.go
```

## Future Improvements

- **Extended Comment Parsing:** Enhance the comment detection logic to support more programming languages and comment styles.
- **Concurrency:** Implement concurrent file processing to improve performance in large projects.
- **Visualization:** Improved visualization in terminal and possible usage of html for in browser inspection

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


