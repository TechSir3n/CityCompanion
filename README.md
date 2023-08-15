# CityCompanion Telegram Bot

CityCompanion is a Telegram bot written in Golang that provides users with various features to explore their city. The bot allows users to search for nearby places based on their current location, choose specific categories of interest, save places, add them to favorites, leave reviews, view reviews, see places on a map, and get walking and driving distances to a particular place.

## Features

- Search for nearby places based on user's location
- Choose categories of places to search
- Save places to view later
- Add places to favorites
- Leave and view reviews for places
- View places on a map
- Get walking and driving distances to a place
- Customize the number of places to view
- Customize the number of photos to view for each place
- Adjust the search radius

## Technologies Used

- Golang
- Telegram Bot API
- Foursquare Maps API

## Installation

1. Clone the repository:

git clone https://github.com/TechSir3n/CityCompanion.git

2. Install the required dependencies:

go mod tidy

3. Set up the required API keys:

- Telegram Bot API: [Create a new bot](https://core.telegram.org/bots#botfather) and obtain the API token.
- Google Maps API: [Create a new project](https://developers.google.com/maps/gmp-get-started#create-project) and enable the Places API. Obtain the API key.

4. Set the API keys as environment variables:

export TELEGRAM_API_TOKEN=your-telegram-api-token
export GOOGLE_MAPS_API_KEY=your-google-maps-api-key

5. Build and run the bot:

go build .
./city-companion-bot

## Usage

1. Start a chat with the CityCompanion bot on Telegram.
2. Share your location with the bot.
3. Use the available commands and buttons to search for places, save them, add them to favorites, leave reviews, view reviews, view places on a map, and get walking/driving distances.

## Contributing

Contributions are welcome! If you have any ideas or suggestions for improving the CityCompanion bot, please create a new issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).