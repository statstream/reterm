# reterm


## Environment Variables
`REDIS_URL`
### Default Value: `"redis://localhost:6379"`

`LOG_FILE`
### Default Value: `reterm.log`

`LOG_MAX_SIZE`
### Default Value: `10`

`LOG_MAX_AGE`
### Default Value: `30`


`LOG_MAX_BACKUPS`
### Default Value: `5`

## Key

- **q**: Quit the application.
- **r**: Refresh the data to update key list and values.
- **d**: Delete the selected key. A confirmation prompt will appear.
- **h**: Display help information.
- **/**: Move focus to the search bar for key search.
## When the search bar is active:
- **Enter**: Search for the provided key.
- **Esc**: Clear the search bar and return focus to the key list.
