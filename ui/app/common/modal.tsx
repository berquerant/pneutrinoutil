import { useCallback, useState } from 'react'
import Modal from 'react-modal'
import { CodeBlock } from './code'

enum CopyState {
  Init,
  Success,
  Failed,
}

function toCopyButtonText(s: CopyState): string {
  switch (s) {
    case CopyState.Init:
      return 'Copy'
    case CopyState.Success:
      return 'Copied!'
    default:
      return 'Copy failed'
  }
}

function toCopyButtonClassName(s: CopyState): string {
  switch (s) {
    case CopyState.Init:
      return 'btn-primary'
    case CopyState.Success:
      return 'btn-success'
    default:
      return 'btn-danger'
  }
}

export default function CodeModal({ name, code }) {
  const [modalOpen, setModalOpen] = useState(false)
  const openModal = () => setModalOpen(true)
  const closeModal = () => setModalOpen(false)

  const [copyState, setCopyState] = useState(CopyState.Init)
  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopyState(CopyState.Success)
      setTimeout(() => {
        setCopyState(CopyState.Init)
      }, 3000)
    } catch (err) {
      setTimeout(() => {
        setCopyState(CopyState.Failed)
      }, 3000)
    }
  }, [code])
  const copyButtonClassName = 'btn ' + toCopyButtonClassName(copyState)
  const copyButtonText = toCopyButtonText(copyState)
  return (
    <div>
    <button type="button" className="btn btn-primary" onClick={openModal} onRequestClose={closeModal}>
    {name}
    </button>
      <Modal isOpen={modalOpen}>
      <div className="row align-items-start">
      <div className="col d-flex gap-3">
      <button type="button" className="btn btn-danger" onClick={closeModal}>Close</button>
      <button type="button" className={copyButtonClassName} onClick={handleCopy}>{copyButtonText}</button>
      </div>
      </div>
      <hr/>
      <div className="row">
      {CodeBlock({code: code})}
      </div>
    </Modal>
    </div>
  )
}
