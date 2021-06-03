{{define "body"}}

{{if .Error}}
    <div class="alert alert-danger">
        <strong>Error!</strong> {{.Error}}
    </div>
{{end}}

{{if .Downloadable}}
    <h1>ファイルのダウンロード</h1>
    <form action="/{{.File.Slug}}/download" method="post">
        <p class="mt20">以下のファイルをダウンロードします。</p>
        <table class="table mt20">
            <tr>
                <th width="200">ファイル名</th>
                <td>{{.File.Name}}</td>
            </tr>
            <tr>
                <th>有効期限</th>
                <td>{{.File.FormatDeadline}}</td>
            </tr>
            {{if ne .File.DownloadLimit 0}}
                <tr>
                    <th>残りダウンロード可能回数</th>
                    <td>{{.File.RemainingDownloadableCount}}回</td>
                </tr>
            {{end}}
        </table>
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
