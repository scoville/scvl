{{define "body"}}

{{if .Success}}
    <div class="alert alert-success">
        <strong>Success!</strong> {{.Success}}
    </div>
{{end}}

<h1>ファイルの編集</h1>
<form action="/{{.File.Slug}}" method="post" enctype="multipart/form-data">
    <div class="form-group">
        <label for="download_limit">ダウンロード制限回数（無制限の場合は0を指定）</label>
        <input type="number" name="download_limit" value="{{.File.DownloadLimit}}" class="form-control" required />
    </div>
    <div class="form-group">
        <label for="password">パスワード（変更しない場合は空白）</label>
        <input type="password" name="password" value="" class="form-control" />
    </div>
    <input type="file" name="file" required>
    <div class="mt20">
        <input type="submit" value="更新" class="btn btn-primary">
    </div>
</form>

<a class="btn btn-default mt20" href="/">ファイル一覧に戻る</a>

{{end}}
