# gal - AUTHORSファイルをgit logから生成
galコマンドは、カレントディレクトリにAUTHORS.mdを生成します。galコマンドは、著者名とメールアドレスをgit logの情報から取得します。これらの情報は、アルファベット順でAUTHORS.mdに記載されます。

# インストール方法
## Step.1 golangのインストール
Golangをシステムにインストールしていない場合は、まずはgolangをインストールしてください。インストール方法は、[Go公式サイト](https://go.dev/doc/install) で確認してください。  
## Step2. galのインストール
```
$ go install github.com/nao1215/gal/cmd/gal@latest
```

# 使い方
".git"ディレクトリが存在するディレクトリで、galコマンドを実行してください。**既存のAUTHORS.mdは上書きされるため、注意してください。**
You execute the gal command in the directory where .git exists. **Please note that the existing AUTHORS.md will be overwritten.**
```
$ gal

$ cat AUTHORS.md 
# Authors List (in alphabetical order)
CHIKAMATSU Naohiro<n.chika156@gmail.com>
TEST User<test@gmail.com>
```

# 連絡先
「バグを見つけた場合」や「機能追加要望」に関するコメントを開発者に送りたい場合は、以下の連絡先を使用してください。

- [GitHub Issue](https://github.com/nao1215/gal/issues)

# ライセンス
galプロジェクトは、[Apache License 2.0](./LICENSE)条文の下でライセンスされています。