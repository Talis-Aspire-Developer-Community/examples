# Example - Iterating a csv of list ids to retrieve the list titles

An example of a simple golang command line tool which will get a list via the RL API and print out it's title to the screen.

You will need to set two envvars in terminal, specifying your API credentials:

```bash
export ACTIVE_TALIS_PERSONA_ID="my-id"
export ACTIVE_TALIS_PERSONA_SECRET="my-secret"
```

Then you can run the tool like this:

```bash
go run ./ -tenant broadminster -infile input.csv -outfile ./test/out.csv
```

If you are trying to use this example for a Canadian tenancy you will need to modify both the persona token URL and the RL API Base URL in the code,
as this example presumes the use of the non-canadian urls.
