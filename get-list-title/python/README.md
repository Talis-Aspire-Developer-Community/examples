# Example - getting a list title in python

An example of a simple python command line tool which will get a list via the Talis Aspire Reading Lists API and print out it's title to the screen.

## Getting setup

We used a MacOS machine when writing this guide, Exact commands may vary on your machine, but the principle is the same.

1. First make sure you have environment variables set.
2. Then setup your python environment

### 1. Environment variables

You will need to set two environment variables in your terminal, specifying your API credentials:

```bash
export ACTIVE_TALIS_PERSONA_ID="my-id"
export ACTIVE_TALIS_PERSONA_SECRET="my-secret"
```

### 2. Virtual Python Environment

This tool uses a [pipenv virtual environment](https://pipenv.pypa.io/en/latest/) to load all the required libraries for you.  This is a recommended practice for all python projects as it helps when you need specific dependencies.

The following command will create the python virtual environment using the `Pipfile` and `Pipfile.lock` to ensure it is setup to the same as we used when writing the script.

```bash
# to crete a new pipenv
# cd into the directory where you extracted this repository.
pipenv shell
# you will now be in a new shell which is using the specific python environment.
```

## Running the tool

Then you can run the tool like this (obviously you will need to use your own tenant short code and example list):

```bash
python get_list_title.py -tenant broadminster -id 06C17A4C-4D2C-D974-A35A-E7C31EC70725 
```

:canada: If you are trying to use this example for a Canadian tenancy you will need to modify both the persona token URL and the RL API Base URL in the code, as this example presumes the use of the non-canadian urls. Search for 'canada' in the code to find where to change these.
