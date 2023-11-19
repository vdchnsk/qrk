### i-go interpreter

Parser implements [**Pratt parsing algorithm**](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/)

### Code snippets

```haskell
let fibanacci = fn(n) {
    if (n < 2) {
        return n;
    }
    return fibanacci(n-2) + fibanacci(n-1)
}
fibanacci(42)
```

```haskell
let isAnswerOnTheLifeQuestion = fn(answer) { 
    let expectedAnswer = 42;
    return answer == expectedAnswer;
};
isAnswerOnTheLifeQuestion(69);
```

### 🚀 How to run locally

- have **go** installed locally

- install dependecies

```bash
go mod download
```

- launch REPL

```bash
make run
```
![image](https://github.com/vdchnsk/i-go/assets/64404596/e71e4366-4ea4-4e6a-8a09-8464bd16db43)
