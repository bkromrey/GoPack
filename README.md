# GoPack
A terminal based packing list application written using Go and bubbletea. 

# Setup & Run
Clone the repo and create a `.env` file in the root directory. 

### Sample .env File
```
# Generate a free API key from Geocoding API for geocode lookups based on Open Street Map.
# Get the API key at https://geocode.maps.co/ and note that the key must be formatted as 
# a string. For example, GEO_API_KEY="123456789" works, while GEO_API_KEY=123456789 wouldn't.
GEO_API_KEY=
```

Then use `go run main.go` to run the program. 