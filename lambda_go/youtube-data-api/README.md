# youtube-data-api

これはyoutube-data-apiのサンプルテンプレートです - 以下、生成したものを簡単に説明します：

```bash
.
├── README.md                   <-- This instructions file
├── hello-world                 <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
│   └── Dockerfile              <-- Dockerfile
└── template.yaml
```

## 必要条件

* AWS CLIはAdministrator権限で設定済みです。
* [Docker installed](https://www.docker.com/community-edition)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

ローカルテストには、以下が必要になる場合があります。
* [Golang](https://golang.org)

## セットアップの流れ

### 依存関係のインストールとターゲットのビルド

この例では、組み込みの `sam build` を使って Dockerfile から docker イメージを構築し、Docker イメージの中にアプリケーションのソースをコピーしています。
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

### ローカル開発環境

**ローカルAPI Gatewayを利用してローカルに関数を呼び出す。**

```bash
sam local start-api
```

先ほどのコマンドが正常に実行されれば、次のようなローカルエンドポイントで関数を呼び出すことができるようになるはずです。
 `http://localhost:3000/hello`

**SAM CLI** は Lambda と API Gateway の両方をローカルでエミュレートするために使用され、この環境の起動方法 (ランタイム、ソースコードの場所など) を理解するために `template.yaml` を使用します - 以下の抜粋は API とそのルートを初期化するために CLI が読むものです：

```yaml
...
Events:
    HelloWorld:
        Type: Api # More info about API Event Source: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
        Properties:
            Path: /hello
            Method: get
```

## パッケージングとデプロイ

AWS Lambda Golang ランタイムは、ビルドステップで生成された実行ファイルを含むフラットフォルダを必要とします。SAM は `CodeUri` プロパティを使用して、アプリケーションを検索する場所を特定します：

```yaml
...
    FirstFunction:
        Type: AWS::Serverless::Function
        Properties:
            CodeUri: hello_world/
            ...
```

アプリケーションを初めてデプロイする場合は、シェルで次のように実行します：

```bash
sam deploy --guided
```

このコマンドは、一連のプロンプトとともに、アプリケーションをパッケージ化してAWSにデプロイします：

* **Stack Name**： スタック名**：CloudFormationにデプロイするスタックの名前です。これはアカウントとリージョンに固有であるべきで、良い出発点はプロジェクト名と一致するものでしょう。
* **Region**： アプリをデプロイするAWSリージョンです。
* **Confirm changes before deploy**： Yesに設定すると、手動で確認するために、実行前に変更セットを表示します。no に設定すると、AWS SAM CLI はアプリケーションの変更を自動的にデプロイします。
* SAM CLI の IAM ロールの作成を許可する**： この例を含む多くの AWS SAM テンプレートは、AWS サービスにアクセスするために含まれる AWS Lambda 関数に必要な AWS IAM ロールを作成します。デフォルトでは、これらは必要最小限の権限にスコープダウンされています。IAM ロールを作成または変更する AWS CloudFormation スタックをデプロイするには、`capabilities` の `CAPABILITY_IAM` 値を提供する必要があります。このプロンプトで権限が提供されない場合、この例をデプロイするには、`--capabilities CAPABILITY_IAM` を `sam deploy` コマンドに明示的に渡す必要があります。
* **Save arguments to samconfig.toml**： yes に設定すると、選択した内容がプロジェクト内の設定ファイルに保存され、将来、アプリケーションに変更を加える際に、パラメータなしで `sam deploy` を再実行できるようになります。

デプロイ後に表示される出力値から、API Gateway EndpointのURLを確認することができます。

### テスト

Golangに組み込まれている `testing` パッケージを使用し、以下のコマンドを実行するだけでローカルにテストを実行することができます：

```shell
go test -v ./hello-world/
```
# Appendix

### Golangのインストール

Go 1.x（xは最新バージョン）がgolang公式サイトの説明に従ってインストールされていることを確認してください： https://golang.org/doc/install

手っ取り早く始めるには、Homebrewやchocolatey、あるいはLinuxのパッケージマネージャを使うのがよいでしょう。

#### Homebrew (Mac)

ターミナルから次のコマンドを実行します：

```shell
brew install golang
```

すでにインストールされている場合は、以下のコマンドを実行し、最新版であることを確認してください：

```shell
brew update
brew upgrade golang
```

<!-- #### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
``` -->

## 次のステージへ導く

ここでは、この全体的なプロセスがどのように機能するかについて、よりよく理解するために使用できるいくつかのアイデアを紹介します：

* 追加のAPIリソース（例：/hello/{proxy+}）を作成し、この新しいパスを通して要求された名前を返す。
* ユニットテストを更新する。
* パッケージ＆デプロイ

次に、hello worldのサンプルの先や、他の人がどのようにServerlessアプリケーションを構成しているのかを知るために、以下のリソースを利用することができます：

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)
