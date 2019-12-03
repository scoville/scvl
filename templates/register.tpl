{{define "body"}}
<div class="jumbotron">
  <h2>ユーザー登録</h2>
  <form action="/register" method="post">
    <p>メールアドレス： {{.Email}}</p>
    <div class="form-inline register">  
      <p>パスワード：<input id="password" type="password" name="password" class="form-control" placeholder="6文字以上"></p>
      <p>名前：　　　<input id="name" type="text" name="name" class="form-control" placeholder="フルネーム"></p>
      <input type="hidden" name="hash" value="{{.Hash}}">
      <p><input type="submit" value="送信" class="btn btn-primary"></p>
    </div>
  </form>
</div>
{{end}}
