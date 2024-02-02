<div align="center">
  <h3 align="center">BOT</h3>
  <p align="center">
    DISCORD BOT!
    
 
    
  </p>
</div>



### Get Started

0. Requirements
    - Weather api key, you can get it free from  [https://www.weatherapi.com/](https://www.weatherapi.com/)
    - Cloud google translate api key - [https://cloud.google.com]
    - Discord Bot token  
1. Clone the repo
   ```sh
   git clone https://github.com/iliyasali2107/discord-bot.git
   ```
2. Go to project direcory
   ```sh
   cd discord-bot
   ```
3. Create `config.yaml` near `example.config.yaml` and Enter:
   ```yaml
    bot_token: bot token
    weather_api_key: weather api key
    weather_api_url: http://api.weatherapi.com/v1
    weather_api_current: current.json
    google_translate_api_key: google translate api key
   ```
    if you don't have keys , you can take it from `example.config.yaml`, i will delete it after you check, but there is no bot_token, because discord refreshes it, when pushing it to repo
    UPD: google also refreshed api key, so you need to provide your own api key

4. You can run:
    ```sh
    make run 
    ```
    or run builded file:
    ```sh
    make run-build
    ```
    without make:
    ```sh
    go run ./cmd/bot/bot.go
    ```
    without make builded:
    ```sh
    ./bot
    ```

5. After bot started, go to your discord server and send to see commands:
    ```discord
    !help
    ```
