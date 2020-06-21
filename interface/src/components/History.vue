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
            :must-sort="true"
            :footer-props="{'items-per-page-options': [10, 20, 30]}"
            show-expand
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
          <template v-slot:item.created_at="{ item }">
            {{ moment(item.created_at).format('HH:mm:ss DD.MM.YYYY') }}
          </template>
          <template v-slot:expanded-item="{ headers, item }">
            <td :colspan="headers.length">{{ item.comment || '-'}}</td>
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
        { text: 'ID', value: 'id', sortable: false },
        { text: 'Action', value: 'action', sortable: true },
        { text: 'Locked By', value: 'locked_by', sortable: true },
        { text: 'Time', value: 'created_at', sortable: true },
        { text: '', value: 'data-table-expand' }
      ],
      historyItems: [],
      loading: true,
      options: {
        sortBy: ['created_at'],
        sortDesc: [true]
      },
      totalItems: 0
    }
  },
  methods: {
    close () {
      this.historyItems = []
      this.$emit('close')
    },
    async getHistory () {
      const { sortBy, sortDesc, page, itemsPerPage } = this.options
      let requestParams = {
        page: page,
        'per-page': itemsPerPage,
        'sort-by': sortBy[0],
        'sort-order': sortDesc[0] ? 'desc' : 'asc'
      }
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
