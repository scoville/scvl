{{define "body"}}

{{if .Error}}
    <div class="alert alert-danger">
        <strong>Error!</strong> {{.Error}}
    </div>
{{end}}

{{if .Downloadable}}
    <h1>ファイルのダウンロード</h1>
    <form action="/files/{{.File.Slug}}/download" method="post">
        {{if ne .File.EncryptedPassword ""}}
            <p>このファイルをダウンロードするためにはパスワードが必要です。</p>
            <div class="form-group">
                <label for="password">パスワード</label>
                <input id="password" type="password" name="password" class="form-control">
            </div>
        {{end}}
        <div>
            <input type="submit" value="ダウンロード" class="btn btn-primary">
        </div>
    </form>
{{end}}

{{end}}
