## Asynchronous Agreement Components

### What is it?

Pure go implementation of asynchronous consensus, including byzantine reliable broadcast,
byzantine agreement and so on.

### Start

```
docker build . -t AAC
docker run AAC 2>&/dev/null  # without logs
docker run AAC               # with logs
```

### Design of code

**Named channel**: Naming channel for all protocols in a tree structure. The child protocol
is connect to parent protocol by `.`.

### LICENSE

Licensed by Apache 2.0.



