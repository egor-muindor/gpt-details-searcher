# GPT Details Searcher

## Overview

GPT Details Searcher is a Go-based project designed to search for specific details (like weights) from web pages using a combination of web scraping, web searching, and AI processing. The project integrates with various services such as Brave Search API, OpenAI, and Colly for web scraping.

## Features

- **Web Search**: Uses Brave Search API to find relevant URLs based on a search query.
- **Web Scraping**: Scrapes the content of the URLs to extract text.
- **AI Processing**: Uses OpenAI to process the scraped text and find specific details like weights.

## Prerequisites

- Go 1.16 or later
- A Brave Search API token
- An OpenAI API token
- A `.env` file with the following environment variables:
    - `BRAVE_TOKEN`
    - `OPENAI_ENDPOINT`
    - `OPENAI_TOKEN`
    - `AI_MODEL`

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/egor-muindor/gpt-details-searcher.git
   cd gpt-details-searcher
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Create a `.env` file in the root directory and add your environment variables:
   ```sh
   BRAVE_TOKEN=your_brave_token
   OPENAI_ENDPOINT=your_openai_endpoint
   OPENAI_TOKEN=your_openai_token
   AI_MODEL=your_ai_model
   ```

## Usage

1. Run the main service:
   ```sh
   go run cmd/service.go
   ```

2. The service will perform the following steps:
    - Load environment variables.
    - Initialize the web searcher, scrapper, and AI processor services.
    - Perform a web search based on a predefined query.
    - Scrape the URLs obtained from the web search.
    - Process the scraped text to find specific details like weights.

## Project Structure

- `cmd/service.go`: The main entry point of the application.
- `internal/services/aiprocessor`: Contains the AI processor service for processing text using OpenAI.
- `internal/services/scrapper`: Contains the scrapper service for web scraping using Colly.
- `internal/services/websearcher`: Contains the web searcher service for searching the web using Brave Search API.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.