# gitstore

Simple git key value store.

Stores raw data in files on a branch other than the one you checked out.

## Usage

Writing the value `bar` to a file named `foo`:

```bash
./store write foo bar
```

Then you can read the file with:

```bash
./store read foo
bar
```

Everything will be stored in a commit chain in the `store` branch with customization to come.
