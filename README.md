# Discord Bot with Weather and Translation Features

This Discord bot provides two unique functionalities: weather information and text translation. You can access these features by using the bot's prefix "!" followed by the respective command.

The first command from which you can start from is the 

```bash
!help
```

It will show you all the available functionality of the bot.

## Weather Information

To access weather information, use the following command:

```bash
!weather help
```

This command will provide you with a list of available weather-related commands and options. You can retrieve weather data for specific locations and get details such as temperature, weather description, and humidity.

```bash
!weather get [location] - gets the weather for a specific location
```

Example Screenshot:
<img width="449" alt="Screenshot 2024-02-02 at 19 39 38" src="https://github.com/zhakhangers/go-discord-bot/assets/73777397/e0dc02ec-6be5-4cf8-b5d8-86727ed105c6">


## Text Translation

For text translation, use the following command:

```bash
!translate help
```

This command will display available translation commands and instructions. You can translate text to various languages by specifying the target language code.

```bash
!translate to [2-letters language code] [content] - translates the text after "to" to the given language
```

## Getting Started

To use this bot, you need to set it up by following these steps:

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/zhakhangers/go-discord-bot.git
   ```

Create a .env file in the project directory and add the required API tokens:
 ```bash
BOT_TOKEN=YourDiscordBotToken
GOOGLE_TRANSLATE_TOKEN=YourGoogleTranslateToken
OPENWEATHER_TOKEN=YourOpenWeatherMapToken
```

BOT_TOKEN: Your Discord bot token.
GOOGLE_TRANSLATE_TOKEN: Your Google Cloud Translation API key.
OPENWEATHER_TOKEN: Your OpenWeatherMap API key.
Customize the bot's behavior and features according to your preferences by modifying the code.

Run the bot:
 ```bash
go run main.go
```

Invite the bot to your Discord server using the **OAuth2 URL** generated by Discord.

Start using the bot by sending commands in your server's channels.

With this Discord bot, you can easily access weather information and translate text to different languages. Enjoy exploring its features!

