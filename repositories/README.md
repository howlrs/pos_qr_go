# Repositories by レイヤードアーキテクチャ

## 概要
このパッケージには、POS QRシステムのリポジトリ層が含まれています。データアクセス層として、モデルとデータベース間のインターフェースを提供します。

## テスト実装状況

### ✅ テスト完了済みリポジトリ
すべてのリポジトリについてテスト実装が完了しており、テスト結果はすべて成功しています。

| リポジトリ | テストファイル       | ステータス   |
| ---------- | -------------------- | ------------ |
| Manager    | `manager_test.go`    | ✅ 完了・成功 |
| Seat       | `seat_test.go`       | ✅ 完了・成功 |
| Session    | `session_test.go`    | ✅ 完了・成功 |
| Store      | `store_test.go`      | ✅ 完了・成功 |

## チーム開発規範

### テスト実装について
- **必須**: 新しいリポジトリを追加する際は、対応するテストファイル（`*_test.go`）を必ず作成してください
- **品質保証**: すべてのテストが成功することを確認してからコミットしてください
- **カバレッジ**: リポジトリの主要な機能とエッジケースをカバーするテストを実装してください
- **命名規則**: テストファイルは `{repository_name}_test.go` の形式で命名してください

### テスト実行方法
```bash
# 全テスト実行
go test ./...

# 特定のリポジトリのテスト実行
go test -v {repository_name}_test.go {repository_name}.go

# カバレッジ付きテスト実行
go test -cover ./...
```

### 継続的品質管理
- 新機能追加時は対応するテストケースも追加してください
- リファクタリング後は既存テストがすべて通ることを確認してください
- テストの保守性を考慮し、適切なテストデータとモックを使用してください

## リポジトリ構成

### 共通インターフェース
すべてのリポジトリは `Repository[T]` インターフェースを実装しています：

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    Read(ctx context.Context) ([]*T, error)
    FindByID(ctx context.Context, id string) (*T, error)
    FindByField(ctx context.Context, field string, value interface{}) ([]*T, error)
    UpdateByID(ctx context.Context, id string, updates map[string]interface{}) error
    DeleteByID(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
    Exists(ctx context.Context, id string) (bool, error)
}
```

### モック実装
各リポジトリには対応するモック実装が提供されており、テスト時に使用されます：
- `MockManagerRepository`
- `MockSeatRepository`
- `MockSessionRepository`
- `MockStoreRepository`

### テストカテゴリ

#### 基本機能テスト
- **CRUD操作**: Create, Read, FindByID, UpdateByID, DeleteByID
- **検索機能**: FindByField, Count, Exists
- **エラーハンドリング**: 各操作でのエラー処理

#### ビジネスロジックテスト
- **Manager**: 重複メール検証、パスワード更新、検索条件
- **Session**: ストア別・席別・ステータス別検索
- **Store**: 重複メール検証、パスワードハッシュ化、複数条件検索

#### エッジケーステスト
- **データ整合性**: 存在しないIDでの操作
- **境界値**: nil値、空文字列での検索
- **パフォーマンス**: 大量データでの操作

#### トランザクションテスト
- **データ一貫性**: 作成→検証→更新→削除の一連の流れ
- **同時実行**: 並行読み書き操作
- **エラー回復**: コンテキストキャンセル、タイムアウト

### 実装上の注意点
- すべてのリポジトリはコンテキストベースの操作をサポートしています
- モック実装は `testify/mock` を使用してテスト可能です
- エラーハンドリングは適切に実装され、テストでカバーされています
- データベースクライアントが `nil` の場合、自動的にモック実装が使用されます