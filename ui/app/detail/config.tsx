import type { Route } from "./+types/detail"
import { CodeBlock } from '../common/code'
import { useCallback, useState } from 'react'
import Modal from 'react-modal'

export default function Config(loaderData: Route.ComponentProps) {
  const [modalOpen, setModalOpen] = useState(false)
  const openModal = () => setModalOpen(true)
  const closeModal = () => setModalOpen(false)
  const data = JSON.stringify(loaderData, null, "  ")

  const [isCopied, setIsCopied] = useState(false)
  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(data)
      setIsCopied(true)
      setTimeout(() => {
        setIsCopied(false);
      }, 3000)
    } catch (err) {
      console.error('Failed to copy', err);
    }
  }, [data])
  const copyButtonClassName = "btn btn-primary " + (isCopied ? "btn-success" : "btn-primary")
  const copyButtonText = isCopied ? "Copied!" : "Copy"
  return (
    <div className="container">
    <button type="button" className="btn btn-primary" onClick={openModal} onRequestClose={closeModal}>
    Config
    </button>
      <Modal isOpen={modalOpen}>
      <div class="row align-items-start">
      <div class="col d-flex gap-3">
      <button type="button" className="btn btn-danger" onClick={closeModal}>Close</button>
      <button type="button" className={copyButtonClassName} onClick={handleCopy}>{copyButtonText}</button>
      </div>
      </div>
      <hr/>
      <div class="row">
      {CodeBlock({code: data})}
      </div>
    </Modal>
    </div>
  )
}
