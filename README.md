# go-protoactor-sample

protoactor-go を使った最小のサンプル集です。Ping/Pong と persistence actor の例を含みます。

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

persistence actor サンプルは以下で実行します。
```sh
make run-persistence
```
```
add +1 => 1
add +2 => 3
add +3 => 6
current => 6
current => 6
```

## コマンド
```sh
make build           # ビルド（bin/go-protoactor-sample を生成）
make run             # Ping/Pong サンプルを実行
make run-persistence # persistence actor サンプルを実行
make test            # テスト
make fmt             # gofmt
make tidy            # 依存整理
```

## 構成
- `cmd/main.go` Ping/Pong サンプルのエントリポイント
- `cmd/persistence/main.go` persistence actor サンプルのエントリポイント
- `internal/domain` ドメイン層
- `internal/usecase` ユースケース層
- `internal/interface_adaptor` インターフェースアダプタ層
- `go.mod` / `go.sum` モジュール定義と依存
- `Makefile` 開発用コマンド
- `AGENTS.md` コントリビュータ向けガイド

## 次の一歩
- アクターメッセージを `struct` にして型安全にする
- アクターを分割してルーティングを試す
