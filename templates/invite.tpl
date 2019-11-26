{{define "body"}}
<div class="jumbotron">
  <h1>ユーザーを招待する</h1>
  <form action="/invite" method="post">
    <div class="form-inline">
      <input id="email" type="email" name="email" class="form-control">
      <input type="hidden" value="{{.FromUserID}}">
      <input type="submit" value="送信" class="btn btn-primary">
    </div>
  </form>
</div>
{{end}}
