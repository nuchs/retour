# Retour

Retour ties into the shell history system, storing commands as they are entered into a sqlite database. It allows the user to easily search through their command history.

My motivation for doing this is that I want to be able to search through my entire command history quickly and easily. This is for personal use only, it will only be used on a single computer.

## What gets stored

1. The command itself
2. The timestamp
3. The working directory
4. The exit status
5. The arguments themselves

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
it should have a single table called history
The table should have the following columns
- id
- command
- timestamp
- working directory
- exit status
- arguments
It should have indexs on the columns command, timestamp and working directory

There should be a .config file in the $HOME/.config/retour/ which contains 
- connection string
- retention period
- exclusion list (a set of regexes)
- limit
The format for the config file is toml
Thre should be a config reader which parses the config file and command line options into a config struct.
there should be the following command line options
- -q --query followed by the sql query
- -r --result followed by one of success, failed, all (defaults to all)
- -c --config followed by the config file location
- -t --time-range followed by one of 'today', 'yesterday', 'thelastweek', 'alltime' (defaults to no time period)
- -l --limit followed by an integer (defaults to 100)
- -w --working-directory followed by a string (defaults to no working directory)


There should be a file to interface with the db
- It should be able to create the database if it doesn't exist
- It should be able to create the table if it doesn't exist
- It should be able to insert a record
- It should be able to query the database with a user provided query
- It should be able to query the db with some precanned queries which are parametrized by
  - a time period
  - a result filter
  - the working directory
  - the number of results to return
- Queries should return results in descending order by timestamp
- It should hide all the implemnation details of the db from the consumer only exposing the the functionality required by the application


The fuzzy matcher should be implemented in the main package in a seaprate file (filter.go), it should be independent of the database and the display mode.
- a fuzzy matcher should be initialised with a set of records
- it should have an intially empty filter
- it should expose the current filter and the filtered set of records
- it should have a method to update the filter
- the initial set of matches should be the entire set of records
- The ui should query the fuzzy matcher for to get the current filter and the filtered set of records rather than storing it as state
- the first iteration of the filter should just check if each record contains the filter string in the command or the arguments

There should be a file for the ui
- it should be a tui
- it shoudl use bubbletea
- it should take a set of results (the model)
- it should display the results with the most recent results at the bottom and the oldest results at the top
- there should be a text input control at the bottom of the screen which can be used to filter the results
- when i type the characters should be added to the text input and the models filfer should be updated.
- if the set of records in the model changes the view should be updated
- it should be able to navigate through the results using the arrow keys or ctrl-n ctrl-p
- it should be able to select a result with return
- it should be able to print the selected result to the terminal

This should all be written in a testable way in golang