# go-protoactor-sample

protoactor-go を使った最小の Ping/Pong サンプルです。アクターが Ping を受け取り、Pong を返します。

## 必要条件
- Go 1.25.6 以上（`go.mod` の `go` バージョンに合わせる）

## 使い方
```sh
make run
```

実行すると以下のような出力になります。
```
received: pong
```

## コマンド
```sh
make build  # ビルド（bin/go-protoactor-sample を生成）
make run    # 実行
make test   # テスト
make fmt    # gofmt
make tidy   # 依存整理
```

## 構成
- `cmd/main.go` サンプルのエントリポイント
- `internal/domain` ドメイン層
- `internal/usecase` ユースケース層
- `internal/interface_adaptor` インターフェースアダプタ層
- `go.mod` / `go.sum` モジュール定義と依存
- `Makefile` 開発用コマンド
- `AGENTS.md` コントリビュータ向けガイド

## 次の一歩
- アクターメッセージを `struct` にして型安全にする
- アクターを分割してルーティングを試す
