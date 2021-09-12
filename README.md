# yktr

[esa](https://esa.io/)の記事を展開してブログのように表示するツールです。

## 設定ファイル

デフォルトの設定ファイルは `yktr.toml` です。

```toml
addr = "127.0.0.1" # listenするアドレス（デフォルト `127.0.0.1`）
port = 8080 # listenするポート（デフォルト `8080`）
team = "docs" # esaのチーム名
token = "<YOUR ACCESS TOKEN>" # アクセストークン cf. https://docs.esa.io/posts/102#%E8%AA%8D%E8%A8%BC%E3%81%A8%E8%AA%8D%E5%8F%AF
```

## 使い方

```sh
yktr # -c any-conf.toml
```
