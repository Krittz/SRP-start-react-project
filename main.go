package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::")
	fmt.Println(":::                                                :::")
	fmt.Println(":::              ███████╗██████╗ ██████╗           :::")
	fmt.Println(":::              ██╔════╝██╔══██╗██╔══██╗          :::")
	fmt.Println(":::              ███████╗██████╔╝██████╔╝          :::")
	fmt.Println(":::              ╚════██║██╔══██╗██╔═══╝           :::")
	fmt.Println(":::              ███████║██║  ██║██║               :::")
	fmt.Println(":::              ╚══════╝╚═╝  ╚═╝╚═╝               :::")
	fmt.Println(":::                                                :::")
	fmt.Println("::::::::::::::::::::::::::::::::::::::::::::::::::::::")
	fmt.Println("")
	fmt.Println("★ Título:        Start React Project")
	fmt.Println("★ Autor:         Cristian Alves Silva")
	fmt.Println("★ Empresa:       Kodificar | Do Brasil para o Mundo")

	// Entrada do nome do projeto
	fmt.Print("\n➤ ➤ ➤  Digite o nome do projeto: ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	// Criação do scaffold base com Vite
	fmt.Println("➤ ➤ ➤  Criando projeto com Vite + React + TypeScript...")
	createViteProject(projectName)

	// Entra na pasta do projeto
	projectPath, _ := filepath.Abs(projectName)
	if err := os.Chdir(projectPath); err != nil {
		fmt.Println("🚨 Erro ao entrar na pasta do projeto: ", err)
		return
	}

	// Remove arquivos desnecessários
	fmt.Println("➤ ➤ ➤  Removendo arquivos padrão inúteis...")
	removeUselessFiles()

	// Instala pacotes essenciais
	fmt.Println("➤ ➤ ➤  Instalando dependências base...")
	run("npm", "install")
	run("npm", "install", "tailwindcss", "@tailwindcss/vite", "axios", "react-router-dom", "lucide-react")
	run("npm", "install", "-D", "@types/node")

	// Configura Tailwind no CSS
	fmt.Println("➤ ➤ ➤  Configurando index.css...")
	writeFile("src/index.css", `@import "tailwindcss";`)

	// Configura tsconfig para alias @/*
	fmt.Println("➤ ➤ ➤  Editando tsconfig.json...")
	updateTSConfig("tsconfig.json")
	updateTSConfigApp("tsconfig.app.json")

	// Configura vite.config.ts com plugins e alias
	fmt.Println("➤ ➤ ➤  Atualizando vite.config.ts...")
	writeFile("vite.config.ts", viteConfigContent())

	// Inicializa o ShadCN com criação do botão e arquivo de config
	fmt.Println("➤ ➤ ➤  Configurando ShadCN com botão padrão...")
	addDefaultShadCNComponent()

	fmt.Println("➤ ➤ ➤  Reescrevendo arquivos de start...")
	personalizeDefaults()

	// Abre o VSCode
	fmt.Println("➤ ➤ ➤  Abrindo no VS Code...")
	run("code", ".")

	fmt.Println("┍━━ ✅  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┑")
	fmt.Println("      Projeto criado com sucesso!")
	fmt.Println("┕━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ✅  ━━━━━━━┙")

}

// Executa comandos no shell
func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Printf("🚨 Erro ao executar %s %s: %v\n", name, strings.Join(args, " "), err)
		os.Exit(1)
	}
}

// Cria projeto Vite + React + TS com scaffold padrão
func createViteProject(projectName string) {
	fmt.Println("┍━━ 👉 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┑")
	fmt.Println("  Executando: npx create-vite@latest", projectName, "--template react-ts")
	fmt.Println("┕━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 👉 ━━━━━━━┙")
	cmd := exec.Command("npx", "create-vite@latest", projectName, "--template", "react-ts")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("🚨 Erro ao criar projeto Vite:", err)
		os.Exit(1)
	}
}

// Remove arquivos desnecessários do projeto padrão Vite
func removeUselessFiles() {
	files := []string{
		"public/vite.svg",
		"src/assets/react.svg",
		"src/App.css",
	}
	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			fmt.Println("⚠️ Não foi possível remover", file+":", err)
		} else {
			fmt.Println("⛔ Removido:", file)
		}
	}
}

// Cria arquivo com conteúdo fornecido
func writeFile(path string, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Println("🚨 Erro ao escrever o arquivo:", path, err)
		os.Exit(1)
	}
}

// Atualiza tsconfig.json principal com baseUrl e paths
func updateTSConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("⚠️ Não foi possível ler tsconfig.json:", err)
		return
	}
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("🚨 Erro ao parsear tsconfig.json:", err)
		return
	}
	config["compilerOptions"] = map[string]interface{}{
		"baseUrl": ".",
		"paths": map[string][]string{
			"@/*": {"./src/*"},
		},
	}
	config["references"] = []map[string]string{
		{"path": "./tsconfig.app.json"},
		{"path": "./tsconfig.node.json"},
	}
	newData, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(path, newData, 0644)
}

// Atualiza tsconfig.app.json com suporte a paths se necessário
func updateTSConfigApp(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("⚠️ Não foi possível ler", path+":", err)
		return
	}
	content := string(data)
	if !strings.Contains(content, `"paths"`) {
		content = strings.Replace(content, `"compilerOptions": {`,
			`"compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    },`, 1)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			fmt.Println("🚨 Erro ao escrever", path+":", err)
		} else {
			fmt.Println("✅ tsconfig.app.json atualizado com alias @/*")
		}
	} else {
		fmt.Println("ℹ️ tsconfig.app.json já possui paths configurado.")
	}
}

// Cria vite.config.ts padrão com Tailwind e aliases
func viteConfigContent() string {
	return `import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
})
`
}

// Adiciona o componente "button" via ShadCN e cria components.json
func addDefaultShadCNComponent() {
	// cmd := exec.Command("npx", "shadcn@latest", "add", "button")
	cmd := exec.Command("npx", "shadcn@latest", "add", "button", "--yes")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("🚨 Erro ao adicionar componente padrão:", err)
		os.Exit(1)
	}
}

func personalizeDefaults() {
	// Substitui o App.tsx padrão
	appContent := `import { Button } from "./components/ui/button";

function App() {
  return (
    <>
      <div className="min-h-svh flex justify-center items-center">
        <Button>Olá mundo</Button>
      </div>
    </>
  );
}

export default App;
`
	err := os.WriteFile("src/App.tsx", []byte(appContent), 0644)
	if err != nil {
		fmt.Println("🚨 Erro ao personalizar App.tsx:", err)
	} else {
		fmt.Println("✅ App.tsx personalizado com layout padrão Kodificar")
	}

	// Substitui o index.html padrão
	htmlContent := `<!doctype html>
<html lang="pt-BR" class="dark">

<head>
  <meta charset="UTF-8" />
  <link rel="icon" type="image/svg+xml" href="#" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Kodificar | Default</title>
</head>

<body>
  <div id="root"></div>
  <script type="module" src="/src/main.tsx"></script>
</body>

</html>
`
	err = os.WriteFile("index.html", []byte(htmlContent), 0644)
	if err != nil {
		fmt.Println("🚨 Erro ao personalizar index.html:", err)
	} else {
		fmt.Println("✅ index.html configurado com base dark e título Kodificar")
	}
}
