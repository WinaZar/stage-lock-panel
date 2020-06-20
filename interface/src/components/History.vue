<template>
  <v-row justify="center">
    <v-dialog v-model="active" fullscreen hide-overlay transition="dialog-bottom-transition">
      <v-card>
        <v-toolbar dark color="indigo">
          <v-btn icon dark @click="close">
            <v-icon>mdi-close</v-icon>
          </v-btn>
          <v-toolbar-title>History of <span class="font-weight-bold">{{stage}}</span></v-toolbar-title>
        </v-toolbar>
        <v-data-table
            :headers="headers"
            :items="historyItems"
            :options.sync="options"
            :server-items-length="totalItems"
            :loading="loading"
            class="elevation-1"
        >
          <template v-slot:item.locked_by="{ item }">
            {{ item.locked_by || '-' }}
          </template>
          <template v-slot:item.action="{ item }">
            <v-chip
              :color="item.action == 'lock' ? 'red': 'green'"
              :label="true"
            >
              {{item.action}}
            </v-chip>
          </template>
          <template v-slot:item.CreatedAt="{ item }">
            {{ moment(item.CreatedAt).format('HH:mm:ss DD.MM.YYYY') }}
          </template>
        </v-data-table>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import moment from 'moment'
export default {
  name: 'History',
  props: {
    stage: String,
    active: Boolean
  },
  data () {
    return {
      moment: moment,
      headers: [
        { text: 'ID', value: 'ID', sortable: false },
        { text: 'Action', value: 'action', sortable: false },
        { text: 'Locked By', value: 'locked_by', sortable: false },
        { text: 'Time', value: 'CreatedAt', sortable: false }
      ],
      historyItems: [],
      loading: true,
      options: {},
      totalItems: 0
    }
  },
  methods: {
    close () {
      this.historyItems = []
      this.$emit('close')
    },
    async getHistory () {
      const { page, itemsPerPage } = this.options
      let requestParams = { page: page, 'per-page': itemsPerPage }
      try {
        let response = await this.$http.get(`/stages/${this.stage}/history`,
          { params: requestParams }
        )
        this.historyItems = response.data.data.history
        this.totalItems = response.data.data.pagination.total
        this.loading = false
      } catch (e) {
        console.log(e.response.data.message)
      }
    }
  },
  watch: {
    active: {
      async handler (val) {
        if (val === true) {
          // manually trigger changes of options obj
          this.options.page = null
          this.options.page = 1
        }
      }
    },
    options: {
      async handler () {
        await this.getHistory()
      },
      deep: true
    }
  }
}
</script>
