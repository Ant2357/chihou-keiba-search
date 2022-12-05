<script>
  import logo from './assets/images/logo-universal.png'
  import {JsonOutputKeiba} from "../wailsjs/go/main/App.js"

  let resultText = ""
  let url

  function jsonOutputKeiba() {
    JsonOutputKeiba(url).then(result => {
      resultText = result;
    }).catch(error => {
      resultText = error;
    })
  }
</script>

<div class="container-top container vh-100">
  <div class="row vh-100">
    <div class="col align-self-center">
      <div class="card-home card shadow">
        <div class="card-body card-home-text">
          <div class="text-center pt-4">
            <img
              id="logo"
              alt="Wails logo"
              src="{logo}"
              class="card-img-top h-50 w-50"
            >
          </div>

          <div class="container">
            <div class="text-center">
              <h1 class="display-4">こんにちは!</h1>
              <p class="lead">netkeibaのレースURL情報を打ち込むと、ローカルにjsonファイルを出力します</p>
            </div>
          </div>

          <div>
            <div class="mb-3">
              <p class="text-center">{resultText}</p>
            </div>
            <div class="input-group mb-3">
              <input
                type="text"
                class="form-control"
                placeholder="URLを入力"
                bind:value={url}
              >
              <button
                class="btn btn-outline-success"
                type="button"
                on:click={jsonOutputKeiba}
              >
                <i class="fas fa-search"></i> 検索
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
.card-home {
  transform: rotate(2deg);
}
.card-home-text {
  transform:skew(0deg, -2deg);
}

.container-top {
  max-width: 960px;
}
</style>
