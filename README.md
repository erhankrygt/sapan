# SAPAN Trading Strategy

A comprehensive Go application that implements the SAPAN trading strategy with concurrent processing capabilities for both Long and Short scenarios.

## Features

- **Long & Short Scenarios**: Detects both bullish and bearish trading setups
- **Technical Indicators**: EMA, Stochastic RSI, MACD validation
- **Candlestick Patterns**: 1-candlestick Pinbar and 2-candlestick Reversal patterns
- **Concurrent Processing**: Multi-threaded stock analysis with worker pools
- **Real-time Progress**: Live progress tracking during processing
- **Thread-safe Operations**: Safe concurrent access to shared resources
- **Environment Configuration**: Secure configuration via environment variables
- **Flexible API Configuration**: Configurable API URL and endpoints

## Prerequisites

- Go 1.19 or higher
- Alpha Vantage API key (free at https://www.alphavantage.co/support/#api-key)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd sapan
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your actual values
```

## Configuration

Create a `.env` file with the following variables:

```bash
# Required
ALPHA_VANTAGE_API_KEY=your_api_key_here

# Optional (with defaults)
ALPHA_VANTAGE_API_URL=https://www.alphavantage.co/query
WORKER_COUNT=5
REQUEST_DELAY_SECONDS=2
STOCKS_FILE=dist/Stocks.json
OUTPUT_SIZE=200
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ALPHA_VANTAGE_API_KEY` | Yes | - | Your Alpha Vantage API key |
| `ALPHA_VANTAGE_API_URL` | No | https://www.alphavantage.co/query | Alpha Vantage API base URL |
| `WORKER_COUNT` | No | 5 | Number of concurrent workers |
| `REQUEST_DELAY_SECONDS` | No | 2 | Delay between API requests |
| `STOCKS_FILE` | No | dist/Stocks.json | Path to stocks JSON file |
| `OUTPUT_SIZE` | No | 200 | Days of historical data to fetch |

## Usage

### Main Application
```bash
cd cmd/sapan
go run main.go
```

### Demo Application
```bash
cd cmd/sapan
go run demo_concurrent.go
```

### With Custom API URL
```bash
ALPHA_VANTAGE_API_KEY=your_key ALPHA_VANTAGE_API_URL=https://your-proxy.com/query go run main.go
```

## SAPAN Strategy Rules

### Long Scenario (Bullish)
- **EMA Trend**: 20 > 50 > 100 > 200 (uptrend)
- **Stochastic RSI**: K < 30 with bullish crossover
- **MACD**: Bull market OR bear market ≤ 5 candlesticks
- **Patterns**: Long 2-candlestick reversal OR Long pinbar reversal

### Short Scenario (Bearish)
- **EMA Trend**: 20 < 50 < 100 < 200 (downtrend)
- **Stochastic RSI**: K > 70 with bullish crossover
- **MACD**: Bear market OR bull market ≤ 5 candlesticks
- **Patterns**: Short 2-candlestick reversal OR Short pinbar reversal

### Priority System
- Long scenario has priority over Short scenario
- Each stock can only be either Long OR Short (mutually exclusive)

## Project Structure

```
sapan/
├── cmd/sapan/           # Main application entry points
├── internal/
│   ├── config/         # Configuration management
│   ├── data/           # Data fetching and loading
│   ├── indicators/     # Technical indicators (EMA, RSI, MACD)
│   ├── processor/      # Concurrent processing logic
│   ├── strategy/       # SAPAN strategy implementation
│   └── watcher/        # Watch list management
├── models/             # Data models
├── dist/               # Data files
└── .env.example        # Environment variables template
```

## API Rate Limits

The application respects Alpha Vantage API rate limits:
- Free tier: 5 requests per minute, 500 requests per day
- Default configuration: 5 workers with 2-second delays
- Adjust `WORKER_COUNT` and `REQUEST_DELAY_SECONDS` as needed

## Advanced Configuration

### Custom API Endpoints
You can use the `ALPHA_VANTAGE_API_URL` environment variable to:
- Point to a proxy server
- Use a different API endpoint
- Test with mock services

### Example with Proxy
```bash
ALPHA_VANTAGE_API_KEY=your_key ALPHA_VANTAGE_API_URL=https://your-proxy.com/alphavantage go run main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
