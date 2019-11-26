{{define "body"}}
<div class="jumbotron">
  <div>
    <p>新しくユーザーを作成しました。以下のURLを共有し、24時間以内にユーザー登録を完了してください。</p>
    <p>URL： {{.InvitationURL}}</p>
  </div>
  <ul><li><a href="{{.MainHost}}">戻る</a></li></ul>
</div>
{{end}}
