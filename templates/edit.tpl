{{define "body"}}

{{if .Success}}
    <div class="alert alert-success">
        <strong>Success!</strong> {{.Success}}
    </div>
{{end}}
<h1>短縮URLの編集</h1>
<form action="/{{.Page.Slug}}" method="post" class="form-inline">
    <input type="url" name="url" value="{{.Page.URL}}" class="form-control">
    <input type="submit" value="送信" class="btn btn-primary">
</form>
<a class="btn btn-default mt20" href="/">TOPへ戻る</a>

{{end}}
