<script>
  import logo from './assets/images/logo-universal.png'
  import {WebScrapingRace} from "../wailsjs/go/main/App.js"

  let url = ""
  let searchResult = {}
  let message = ""

  function webScrapingRace() {
      message = "ロード中…";
    WebScrapingRace(url).then(result => {
      searchResult = result;
      message = "読み込み完了";
    }).catch(error => {
      message = error;
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
              <p class="lead">netkeibaのレースURL情報を打ち込むと、レース情報が表示されます</p>
            </div>
          </div>

          <div>
            <div class="mb-3">
              <p class="text-center">{message}</p>
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
                on:click={webScrapingRace}
              >
                <i class="fas fa-search"></i> 生成
              </button>
            </div>

            <div class="mb-3">

              {#if Object.keys(searchResult).length !== 0}
                <div class="mb-3">
                  <h4>{searchResult.name}({searchResult.racetrack}{searchResult.distance})</h4>
                </div>

                <table class="table">
                  <thead>
                    <tr>
                      <th scope="col">#</th>
                      <th scope="col">名前</th>
                      <th scope="col">レース数</th>
                      <th scope="col">勝利数</th>
                      <th scope="col">敗北数</th>
                      <th scope="col">コース適正</th>
                      <th scope="col">距離適正</th>
                      <th scope="col">脚質</th>
                      <th scope="col">重馬場</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each searchResult.horses as horse, i}
                      <tr>
                        <th>{i + 1}</th>
                        <th>{horse.name}</th>
                        <th>{horse.play_game_count}</th>
                        <th>{horse.win}</th>
                        <th>{horse.lose}</th>
                        <th>{horse.course_aptitude}</th>
                        <th>{horse.distance_aptitude}</th>
                        <th>{horse.running_style}</th>
                        <th>{horse.heavy_racetrack}</th>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              {/if}
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
