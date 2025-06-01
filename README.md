### Info

boot.dev guided project

this is RSS feed aggregator

built on go version 1.24.1

### Additional tools used

- postgresql is used as db
- goose was used as db migration tool
- sqlc is used to generate go code that interacts with db

### How to use the application

You need ".gatorconfig.json" file in your user profile root directory with database connection information. Example:
`{"db_url":"postgres://db_user:db_pass@localhost:5432/db_schema?sslmode=disable"}`

goose cli tool should be used to set up the database scheema based on the migration files in `/sql/schema` directory

- `gator register \<username\>`                     #registers new user to DB, required for login
- `gator login \<username\>`                        #logs in the user, login info is cached in the .gatorconfig.json
- `gator addfeed \<feed_name\> \<rss_feed_url\>`    #adds new feed to the database and user who does it is automatically marked as the one who follows the feed
- `gator feeds`                                     #lists all feeds in the database
- `gator follow \<rss_feed_url\>`                   #current logged in user is maked to follow the feed based on url if it exists is added to the db
- `gator unfollow \<rss_feed_url\>`                 #current logged in user is removed from following the feed based on url
- `gator agg \<interval_time\>`                     #starts collecting posts from feeds to the db, has to be ended with Ctrl+c. Interval example values: '1m', '10m', '1h'
- `gator browse`                                    #lists most recent posts from feeds the current user is following. Takes optional number argument on how many posts are show, default is 2

### boot.dev recommended extension ideas for this project

- Add sorting and filtering options to the browse command
- Add pagination to the browse command
- Add concurrency to the agg command so that it can fetch more frequently
- Add a search command that allows for fuzzy searching of posts
- Add bookmarking or liking posts
- Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
- Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
- Write a service manager that keeps the agg command running in the background and restarts it if it crashes
