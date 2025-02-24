# Retour

Retour ties into the shell history system, storing commands as they are entered into a sqlite database. It allows the user to easily search through their command history.

My motivation for doing this is that I want to be able to search through my entire command history quickly and easily. This is for personal use only, it will only be used on a single computer.

## What gets stored

1. The command itself
2. The timestamp
3. The working directory
4. The exit status
5. The number of arguments
6. The arguments themselves

## Querying the database
We want to be able to query the database in different ways.

1. An interactive mode which shows `all the data but filters it down as we type characters
  - We can filter the command by typing a string which will be fuzzy matched against the command and the arguments
  - We can use the arrow keys to navigate through the results.
  - Hitting return selects the command and prints it to the terminal.
  - this is the default mode
2. A query mode which allows us to write a full sql query. It prints the results to the terminal.


## Structure
The database is a sqlite database.

There should be a .config file in the $HOME/.config/retour/ which contains 
- connection string
- retention period
- exclusion list (a set of regexes)
The format for the config file is toml
Thre should be a config reader which parses the config file and command line options into a config struct.
there should be the following command line options
- -q --query followed by the sql query
- -r --result followed by one of success, failed, all (defaults to all)
- -c --config followed by the config file location
- -t --time-range followed by one of 'today', 'yesterday', 'thelastweek', 'alltime' (defaults to no time period)


There should be a package to interface with the db
- It should be able to create the database if it doesn't exist
- It should be able to create the table if it doesn't exist
- It should be able to insert a record
- It should hide all the implemnation details of the db from the consumer only exposing the the functionality required by the application

The fuzzy matcher should be implemented as a separate package, it should be independent of the database and the display mode.
- a fuzzy matcher should be initialised with a set of records
- it should have an intially empty filter
- the initial set of matches should be the entore set of records
- it should be possible to add characters to the filter

There should be a package for the display mode
- it should be a tui
- it should be able to display the results of a query
- it should be able to display the results of a fuzzy match
- it should be able to navigate through the results
- it should be able to select a result
- it should be able to print the selected result to the terminal

This should all be written in a testable way in golang