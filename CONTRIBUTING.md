# Contribuindo ao Almodon

Contribuições ao Almodon são limitadas, atualmente, apenas ao integrantes da sua equipe de desenvolvimento, reporte de bugs, críticas e sugestões são bem-vindas através do sistema de issues do GitHub.

## Diretrizes para a Padronização de Código

Essa seção lista todas as condições consideradas essenciais para o cultivo de uma base de código sem surpresas. A violação de qualquer uma dessas diretrizes resultará em rejeição de pull requests.

### 1. Todos os identificadores devem estar em inglês

Como a linguagem Go tem sua estrutura de controle e biblioteca padrão escritas em inglês, é natural que todo o código seja em inglês.

### 2. Padrão de nomeação de identificadores

Identificadores exportados, naturalmente, seguem PascalCase, até mesmo constantes e variáveis exportadas.

Identificadores não constantes não exportados e parametros e variáveis em escopos seguem camelCase.

(Talvez seja revogada) Identificadores constantes não exportados seguem _PascalCase, que é o mesmo que PascalCase, mas com um underline (`_`), prefixado.

Identificadores de rótulos seguem PascalCase.

### 3. Nomeação de funções e métodos, tipos e variáveis

Funções e método são nomeadas a partir do que fazem e retornam, e caso necessário, para suprir a ausência de sobrecarga, o que recebem. É preferível `<ação> [ <objeto de ação> ]`. Para métodos, é comum que o receptor tenha apenas uma ou duas letras, pois seu tipo já apresenta informações suficiente. Também não é incomum para funções e baixo nível terem nomes com apenas uma ou duas letras, porém, para as políticas de alto nível, essa prática é desencorajada.

Para tipos, a menos que esteja referindo-se à uma entidade (no contexto de modelagem), ela não precisa carregar o nome da entidade. Por exemplo, em Java, é comum que exista `UserController`, já em Go, como o pacote é sempre referenciado por todos os seus identificadores exportados, `user.UserController` é exagerado, portando conforma-se apenas a `Controller`.

Para entidades, tipos que são um subconjunto de sua representação real obrigada por _setters_, deve-se seguir o seguinte modelo:

```go
type User struct {
    name string
}

func (u *User) Name() string
func (u *User) SetName(string) error
```

Os _getters_ não possuem "get" no seu nome e os _setters_ sempre possibilitam o retorno de um error, mesmo que, no momento, aceitem todo valor.

Construtores seguem o padrão `New <tipo>`, normalmente retornado um tipo (ou uma referência para esse tipo) e, opcionalmente, um erro. Como:

```go
func NewUser(string) (User, error)
```

Se um pacote compreender apenas um tipo expressivo, como o pacote `errors` da biblioteca padrão, usar apenas `New` é encorajado.

### 4. Código deve ser formatado pelo `go fmt`

A ferramenta `fmt` da cadeia de ferramentas do compilador fornece um estilo bem definido, um formato canônico, que faz esse documento ser mais simples. O comando `go fmt` é integrado na grande maioria de editores de texto juntamente a outras ferramentas da linguagem.

### 5. Escopos devem não apresentar linhas em branco no começo ou no final

Escopos, blocos de código delimitados por `{` e `}`, devem sempre possuir uma declaração no início e final, ou seja, o código a seguir é encorajado:

```go
func main() {
    var num int
    fmt.Print("Enter a number: ")
    fmt.Scanf("%d\n", &num)

    if num < 0 {
        fmt.Printf("The number %d is non-positive")
    } else {
        fmt.Printf("The number %d is positive")
    }
}
```

Enquanto o código a seguir será rejeitado:

```go
func main() {

    var num int
    fmt.Print("Enter a number: ")
    fmt.Scanf("%d\n", &num)

    if num < 0 {
        fmt.Printf("The number %d is non-positive")

    } else {
        fmt.Printf("The number %d is positive")
    }
}
```

### 6. Formato de commits

Com exceção do primeiro commit "ready, set, Go!", todo commit deve seguir a estrutura:

```go
<tipo> [ "(" <escopo> ")" ] [ "!" ] ": " <mensagem>
```

O `<tipo>` pode ser um dos valores a seguir:

- `feat`: para criação de novas funcionalidades
- `doc`: para a documentação de funcionalidades
- `fix`: para o conserto de bugs
- `test`: para a adição de testes

O `<escopo>` refere-se ao pacote ou unidade semântica afetada pelo commit, mudanças que envolvem mais de um pacote ou unidade semântica em um único commit devem ser evitadas.

O `!` opcional descreve `BREAKING CHANGE`, que é presente quando uma API pública é alterada, que pode fazer outros pacotes, em cascata, terem que ser alterados.

A `<mensagem>` é uma descrição curta, normalmente uma única oração, que descreve as alterações feita nesse commit.