<template>
  <div id="app" class="container">
    <div class="section has-text-centered">
      <h1 class="title is-1">Stage servers status panel</h1>
    </div>
    <div class="stage section" v-for="item in stages" :key="item.name">
      <div class="card">
        <header class="card-header">
          <p class="card-header-title">{{item.name}}</p>
        </header>
        <div class="card-content">
          <div class="content">
            <p>Status: <span
                        class="status"
                        :class="[item.locked ? 'has-background-danger': 'has-background-success']"
                        >
                          {{item.locked ? 'Locked': 'Free'}}
                      </span>
            </p>
            <p v-if="item.locked">Locked by <b>{{item.locked_by}}</b></p>
            <p v-if="item.comment.length > 0">Comment: {{item.comment}}</p>
          </div>
        </div>
        <footer class="card-footer">
          <div class="card-footer-item">
            <div class="buttons">
              <b-button @click="lockDialog(item.name)" type="is-success" :disabled="item.locked">Lock</b-button>
              <b-button @click="unlockDialog(item.name)" type="is-danger" :disabled="!item.locked">Unlock</b-button>
            </div>
          </div>
        </footer>
      </div>
    </div>
  </div>
</template>

<script>
import LockDialog from '@/components/LockDialog.vue'
import UnLockDialog from '@/components/UnlockDialog.vue'

export default {
  name: 'App',
  data () {
    return {
      backUrl: process.env.VUE_APP_BACKEND_URL,
      stages: []
    }
  },
  methods: {
    async lockDialog (stageName) {
      this.$buefy.modal.open({
        parent: this,
        component: LockDialog,
        hasModalCard: true,
        trapFocus: true,
        props: {
          stage: stageName
        },
        events: {
          refresh: this.fetchStagesStatuses
        }
      })
    },
    async unlockDialog (stageName) {
      this.$buefy.modal.open({
        parent: this,
        component: UnLockDialog,
        hasModalCard: true,
        trapFocus: true,
        props: {
          stage: stageName
        },
        events: {
          refresh: this.fetchStagesStatuses
        }
      })
    },
    async fetchStagesStatuses () {
      let response = await this.$http.get('/stages')
      this.stages = response.data.data
    }
  },
  async created () {
    await this.fetchStagesStatuses()
  }
}
</script>

<style lang="css">
  .stage.section {
    padding: 1rem 1.5rem;
  }

  span.status {
    padding: 0.5rem;
  }

</style>
