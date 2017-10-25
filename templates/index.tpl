{{define "body"}}
<div class="jumbotron">
  <h1>SCVL URL Shortener</h1>
  <p>URLの短縮ができます。</p>
  <form action="/shorten" method="post" class="form-inline">
    <input type="url" name="url" value="{{.URL}}" class="form-control">
    <input type="submit" value="送信" class="btn btn-primary">
  </form>
  {{if ne .URL ""}}
    <div class="panel panel-success">
      <div class="panel-heading">
        <h3 class="panel-title">短縮が成功しました。</h3>
      </div>
      <div class="panel-body">
        <p>短縮結果:</p>
        <p><script>document.write(location.origin)</script>/{{.Slug}}</p>
      </div>
    </div>
  {{end}}
</div>
{{end}}
