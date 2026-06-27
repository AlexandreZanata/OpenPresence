import { DEV_MOCK_USERS } from '~/lib/auth/dev-users'
import { loginPageStyles as styles } from './form-styles'

type Props = {
  onPick: (registrationId: string, password: string) => void
}

export function MockCredentialsHint({ onPick }: Props) {
  return (
    <div style={styles.mockHint}>
      <p style={styles.mockHintTitle}>Dev test accounts</p>
      <ul style={styles.mockHintList}>
        {DEV_MOCK_USERS.map((user) => (
          <li key={user.registrationId}>
            <button
              type="button"
              style={styles.mockHintButton}
              onClick={() => onPick(user.registrationId, user.password)}
            >
              <strong>{user.registrationId}</strong> / {user.password}
              <span style={styles.mockHintRole}> ({user.roles.join(', ')})</span>
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}
