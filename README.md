# ⚡ TD-Spider ⚡
> Text Density Crawler

This project is a simple web crawler that searches for a keyword from a starting URL and crawls through connected web pages. 

It extracts text from web pages using the techniques from the "***DOM Based Content Extraction via Text Density***" research paper and prints the text containing the specified keyword along with the URL where it was found.

## Benchmark Performance
The following benchmark performance was observed on a MacBook Pro with the following specifications:

```
OS version: 13.2.1(22D68)
Processor: 2.3 GHz Quad-Core Intel Core i5
Memory: 16GB 2133 MHz LPDDR3
Graphics: Intel Iris Plus Graphics 655 1536 MB
```

> The performance was measured by running the crawler for 10 seconds:
```
[ 1 SECOND] REQUEST_COUNT :  276 PARSING_COUNT :  74
[ 2 SECOND] REQUEST_COUNT :  567 PARSING_COUNT :  361
[ 3 SECOND] REQUEST_COUNT :  937 PARSING_COUNT :  731
[ 4 SECOND] REQUEST_COUNT :  1192 PARSING_COUNT :  985
[ 5 SECOND] REQUEST_COUNT :  1525 PARSING_COUNT :  1317
[ 6 SECOND] REQUEST_COUNT :  1811 PARSING_COUNT :  1601
[ 7 SECOND] REQUEST_COUNT :  1993 PARSING_COUNT :  1785
[ 8 SECOND] REQUEST_COUNT :  2194 PARSING_COUNT :  1985
[ 9 SECOND] REQUEST_COUNT :  2373 PARSING_COUNT :  2163
[ 10 SECOND] REQUEST_COUNT :  2503 PARSING_COUNT :  2292
```


## Getting Started

These instructions will help you set up and run the project on your local machine.

### Prerequisites

To run this project, you need to have Go installed on your system. You can download it from the [official website](https://golang.org/dl/).

### Running the Crawler

1. Clone the repository to your local machine:

> git clone https://github.com/LandWhale2/TD-Spider.git


2. Navigate to the project directory:

> cd TD-Spider


3. Run the crawler with a starting URL and the keyword you want to search for:

> go run spider.go <URL> <KEYWORD>



Replace `<URL>` with the starting URL, and `<KEYWORD>` with the keyword you want to search for.

## Contributing

We welcome contributions from everyone! If you'd like to contribute to this project, please feel free to submit a pull request or open an issue. All contributions must follow the standard [GitHub Flow](https://guides.github.com/introduction/flow/).
  
## Acknowledgments

This project was inspired by and based on the following research paper:

- **Title**: DOM Based Content Extraction via Text Density

We would like to express our gratitude to the authors for their valuable insights and guidance in the development of this simple web crawler.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
