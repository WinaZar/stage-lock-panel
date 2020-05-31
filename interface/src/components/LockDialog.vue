<template>
  <div class="modal-card">
    <header class="modal-card-head">
      <p class="modal-card-title">Lock {{stage}}</p>
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
      <b-field
        label="Lock Owner"
        :type="getLableType('lockOwner')"
        :message="getLabelMessage('lockOwner')"
      >
        <b-input
          placeholder="Lock Owner"
          v-model.trim="$v.lockOwner.$model"
          type="text"
          name="lockOwner"
        >
        </b-input>
      </b-field>
      <b-field
        label="Lock Comment"
        :type="getLableType('lockComment')"
        :message="getLabelMessage('lockComment')"
      >
        <b-input
          placeholder="Lock Comment"
          v-model.trim="$v.lockComment.$model"
          type="textarea"
          name="lockComment"
        >
        </b-input>
      </b-field>
    </section>
    <footer class="modal-card-foot">
      <div class="buttons">
        <b-button class="is-success" @click="lock">Lock</b-button>
        <b-button class="is-primary" @click="$parent.close()">Close</b-button>
      </div>
    </footer>
  </div>
</template>

<script>
import { required, minLength, maxLength } from 'vuelidate/lib/validators'
export default {
  name: 'LockDialog',
  props: {
    stage: String
  },
  data () {
    return {
      lockCode: '',
      lockOwner: '',
      lockComment: ''
    }
  },
  validations: {
    lockCode: {
      required,
      minLength: minLength(3),
      maxLength: maxLength(8)
    },
    lockOwner: {
      required,
      minLength: minLength(3),
      maxLength: maxLength(100)
    },
    lockComment: {
      maxLength: maxLength(500)
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
    async lock () {
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
        let response = await this.$http.post(`/stages/${this.stage}/lock`, {
          code: this.lockCode,
          locked_by: this.lockOwner,
          comment: this.lockComment
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
