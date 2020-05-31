<template>
  <div class="modal-card">
    <header class="modal-card-head">
      <p class="modal-card-title">Unlock {{stage}}</p>
    </header>
    <section class="modal-card-body">
      <b-field
        label="Lock Code"
        :type="getLableType('lockCode')"
        :message="getLabelMessage('lockCode')"
      >
        <b-input
          placeholder="Lock Code"
          v-model.trim="$v.lockCode.$model"
          type="text"
          name="lockCode"
        >
        </b-input>
      </b-field>
    </section>
    <footer class="modal-card-foot">
      <div class="buttons">
        <b-button class="is-success" @click="unlock">Unlock</b-button>
        <b-button class="is-primary" @click="$parent.close()">Close</b-button>
      </div>
    </footer>
  </div>
</template>

<script>
import { required, minLength, maxLength } from 'vuelidate/lib/validators'
export default {
  name: 'UnLockDialog',
  props: {
    stage: String
  },
  data () {
    return {
      lockCode: ''
    }
  },
  validations: {
    lockCode: {
      required,
      minLength: minLength(3),
      maxLength: maxLength(8)
    }
  },
  methods: {
    getLableType (name) {
      if (this.$v[name].$error) {
        return 'is-danger'
      }
    },
    getLabelMessage (name) {
      let validation = this.$v[name]
      if (validation.$error) {
        let message = 'Invalid field value'
        if (!validation.required) {
          message = 'Field is required'
        } else if (!validation.minLength) {
          message = `Field must have at least ${validation.$params.minLength.min} letters`
        } else if (!validation.maxLength) {
          message = `Field must have less than ${validation.$params.maxLength.max} letters`
        }
        return message
      }
    },
    async unlock () {
      this.$v.$touch()
      if (this.$v.$invalid) {
        this.$buefy.notification.open({
          message: 'Check form fields',
          type: 'is-danger'
        })
        return false
      }
      let data = null
      try {
        let response = await this.$http.post(`/stages/${this.stage}/unlock`, {
          code: this.lockCode,
          locked_by: 'dummy'
        })
        data = response.data
      } catch (e) {
        data = e.response.data
      }
      this.$buefy.notification.open({
        message: data.message,
        type: data.status === 'error' ? 'is-danger' : 'is-success'
      })
      this.$emit('refresh')
      this.$parent.close()
    }
  }
}
</script>
