# watermill-cli

> [Português Brasileiro](#português-brasileiro) | [English](#english)

---

## English

A CLI for batch video editing — built as a learning project for Go and [Cobra](https://github.com/spf13/cobra). Wraps [ffmpeg-go](https://github.com/u2takey/ffmpeg-go) to automate trimming, concatenation, and adding intros/outros.

### Requirements

- [Go](https://golang.org/dl/) 1.25+
- [FFmpeg](https://ffmpeg.org/download.html) installed and available on PATH

### Build

```sh
git clone https://github.com/seu-usuario/watermill-cli.git
cd watermill-cli
go build -o watermill-cli .
```

To cross-compile for Windows:

```sh
GOOS=windows GOARCH=amd64 go build -o watermill-cli.exe .
```

To compile for macOS (Apple Silicon):

```sh
GOOS=darwin GOARCH=arm64 go build -o watermill-cli .
```

---

### Command 1 — `run`: Automatically edit all videos in a folder

This is the main command. It walks through **all videos** in a folder (and subfolders), trims the beginning/end of each one, and adds intro and outro automatically.

**Expected folder structure:**

```
Videos/
  intro.mp4        ← required (exact name)
  outro.mp4        ← required (exact name)
  lesson1.mp4
  lesson2.mp4
  module2/
    lesson3.mp4
```

**Usage:**

```sh
watermill-cli run Videos/ --removeFirst 3 --removeLast 2
```

**Options:**

| Flag | Default | Description |
|---|---|---|
| `--removeFirst` | `0` | Seconds to remove from the beginning of each video |
| `--removeLast` | `0` | Seconds to remove from the end of each video |
| `--intro` | `intro.mp4` | Path to the intro file |
| `--outro` | `outro.mp4` | Path to the outro file |
| `-v` | — | Show verbose FFmpeg output |

If your intro/outro files have different names:

```sh
watermill-cli run Videos/ --intro opening.mp4 --outro closing.mp4 --removeFirst 3
```

**What happens after running:**

- Each processed video gets an `_edited` suffix (e.g. `lesson1_edited.mp4`)
- A `watermill-cli-progress` file is created in the folder — **do not delete it**. It tracks which videos have already been edited; if you run the command again, only pending videos will be processed

---

### Command 2 — `trim`: Trim the beginning and/or end of a video

Removes seconds from the start or end of a single video.

**Usage:**

```sh
watermill-cli trim lesson.mp4 --removeFirst 3 --removeLast 2 --output lesson_trimmed.mp4
```

**Options:**

| Flag | Default | Description |
|---|---|---|
| `--removeFirst` | `0` | Seconds to remove from the beginning |
| `--removeLast` | `0` | Seconds to remove from the end |
| `-o`, `--output` | `output.mp4` | Output file path |

---

### Command 3 — `concatenate`: Join multiple videos into one

Joins two or more videos in sequence.

**Usage:**

```sh
watermill-cli concatenate part1.mp4 part2.mp4 part3.mp4 --output final.mp4
```

**Options:**

| Flag | Default | Description |
|---|---|---|
| `-o`, `--output` | `output.mp4` | Output file path |
| `-d`, `--dimensions` | `1920:1080` | Output video resolution |

---

## Português Brasileiro

Uma CLI para edição de vídeos em lote — construída como projeto de aprendizado em Go e [Cobra](https://github.com/spf13/cobra). Usa [ffmpeg-go](https://github.com/u2takey/ffmpeg-go) para automatizar cortes, concatenação e adição de intro/outro.

### Dependências

- [Go](https://golang.org/dl/) 1.25+
- [FFmpeg](https://ffmpeg.org/download.html) instalado e disponível no PATH

### Compilação

```sh
git clone https://github.com/seu-usuario/watermill-cli.git
cd watermill-cli
go build -o watermill-cli .
```

Para compilar para Windows a partir de outro sistema:

```sh
GOOS=windows GOARCH=amd64 go build -o watermill-cli.exe .
```

Para compilar para macOS (Apple Silicon):

```sh
GOOS=darwin GOARCH=arm64 go build -o watermill-cli .
```

---

### Comando 1 — `run`: Editar automaticamente todos os vídeos de uma pasta

Esse é o comando principal. Ele percorre **todos os vídeos** de uma pasta (e subpastas), corta o início/fim de cada um, e adiciona intro e outro automaticamente.

**Estrutura de pasta esperada:**

```
Videos/
  intro.mp4        ← obrigatório (nome exato)
  outro.mp4        ← obrigatório (nome exato)
  aula1.mp4
  aula2.mp4
  modulo2/
    aula3.mp4
```

**Sintaxe:**

```sh
watermill-cli run Videos/ --removeFirst 3 --removeLast 2
```

**Opções:**

| Flag | Padrão | Descrição |
|---|---|---|
| `--removeFirst` | `0` | Segundos a remover do início de cada vídeo |
| `--removeLast` | `0` | Segundos a remover do final de cada vídeo |
| `--intro` | `intro.mp4` | Caminho para o arquivo de intro |
| `--outro` | `outro.mp4` | Caminho para o arquivo de outro |
| `-v` | — | Exibe a saída detalhada do FFmpeg |

Se os arquivos de intro/outro tiverem nomes diferentes:

```sh
watermill-cli run Videos/ --intro abertura.mp4 --outro encerramento.mp4 --removeFirst 3
```

**O que acontece após rodar:**

- Cada vídeo processado ganha o sufixo `_edited` (ex: `aula1_edited.mp4`)
- Um arquivo `watermill-cli-progress` é criado na pasta — **não apague ele**. Ele registra quais vídeos já foram editados; se você rodar o comando novamente, apenas os vídeos pendentes serão processados

---

### Comando 2 — `trim`: Cortar o início e/ou fim de um vídeo

Remove segundos do começo ou do final de um único vídeo.

**Sintaxe:**

```sh
watermill-cli trim aula.mp4 --removeFirst 3 --removeLast 2 --output aula_cortada.mp4
```

**Opções:**

| Flag | Padrão | Descrição |
|---|---|---|
| `--removeFirst` | `0` | Segundos a remover do início |
| `--removeLast` | `0` | Segundos a remover do final |
| `-o`, `--output` | `output.mp4` | Caminho do arquivo de saída |

---

### Comando 3 — `concatenate`: Juntar vários vídeos em um só

Une dois ou mais vídeos em sequência.

**Sintaxe:**

```sh
watermill-cli concatenate parte1.mp4 parte2.mp4 parte3.mp4 --output video_final.mp4
```

**Opções:**

| Flag | Padrão | Descrição |
|---|---|---|
| `-o`, `--output` | `output.mp4` | Caminho do arquivo de saída |
| `-d`, `--dimensions` | `1920:1080` | Resolução do vídeo final |
