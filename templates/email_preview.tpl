{{define "body"}}
<div class="email-preview">
  <h1>メールプレビュー</h1>

  <p>以下のメールを{{.EmailTemplate.SendingNumber}}名に送信します。（最初の３件のみ表示）</p>
  <p>メールタイトル</p>
  <p>{{.EmailTemplate.Title}}<p>

  <table class="table">
    <tr>
      <th>送信先</th>
      <td>送信文言</td>
    </tr>
    {{range .EmailTemplate.FilterEmailNum(3)}}
      <tr>
        <td></td>
        <td></td>
      </tr>
    {{end}}
  </table>
  <form action="{{if .IsEmailDomain}}/{{else}}/emails{{end}}" method="post">
    <input type="hidden" name="spreadsheet_url" value="{{.EmailTemplate.BatchEmail.SpreadsheetURL}}" />
    <input type="hidden" name="sheet_name" value="{{.EmailTemplate.BatchEmail.SheetName}}" />
    <input type="hidden" name="sender" value="{{.EmailTemplate.BatchEmail.Sender}}" />
    <input type="hidden" name="title" value="{{.EmailTemplate.BatchEmail.Title}}" />
    <input type="hidden" name="template" value="{{.EmailTemplate.Body}}" />
    <input type="submit" value="送信" class="btn btn-primary">
  </form>
</div>
{{end}}
