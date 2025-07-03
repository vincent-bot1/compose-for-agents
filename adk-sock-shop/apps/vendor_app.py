import streamlit as st
import requests
import json
import uuid
import time

st.set_page_config(
    page_title="Vendor Portal",
    page_icon="🧦",
    layout="wide"
)

API_BASE_URL = "http://adk:8080"
APP_NAME = "agents"

if "user_id" not in st.session_state:
    st.session_state.user_id = f"vendor-{uuid.uuid4()}"

if "session_id" not in st.session_state:
    st.session_state.session_id = None
    
if "messages" not in st.session_state:
    st.session_state.messages = []

def create_adk_session():
    try:
        session_id = f"session={int(time.time())}"
        response = requests.post(
            f"{API_BASE_URL}/apps/{APP_NAME}/users/{st.session_state.user_id}/sessions/{session_id}",
            headers={"Content-Type": "application/json"},
            data=json.dumps({})
        )
        if response.status_code == 200:
            st.session_state.session_id = session_id
            st.session_state.messages = []
            return True
        else:
            st.error(f"Failed to create session: {response.text}")
            return False
    except Exception as e:
        return False

def send_message(message):
    """
    Send a message to the speaker agent and process the response.
    
    This function:
    1. Adds the user message to the chat history
    2. Sends the message to the ADK API
    3. Processes the response to extract text and audio information
    4. Updates the chat history with the assistant's response
    
    Args:
        message (str): The user's message to send to the agent
        
    Returns:
        bool: True if message was sent and processed successfully, False otherwise
    
    API Endpoint:
        POST /run
        
    Response Processing:
        - Parses the ADK event structure to extract text responses
        - Looks for text_to_speech function responses to find audio file paths
        - Adds both text and audio information to the chat history
    """
    if not st.session_state.session_id:
        st.error("No active session. Please create a session first.")
        return False
    
    # Add user message to chat
    st.session_state.messages.append({"role": "user", "content": message})
    
    # Send message to API
    response = requests.post(
        f"{API_BASE_URL}/run",
        headers={"Content-Type": "application/json"},
        data=json.dumps({
            "app_name": APP_NAME,
            "user_id": st.session_state.user_id,
            "session_id": st.session_state.session_id,
            "new_message": {
                "role": "user",
                "parts": [{"text": message}]
            }
        })
    )
    
    if response.status_code != 200:
        st.error(f"Error: {response.text}")
        return False
    
    # Process the response
    events = response.json()
    
    # Extract assistant's text response
    assistant_message = None
    audio_file_path = None
    
    for event in events:
        # Look for the final text response from the model
        if event.get("content", {}).get("role") == "model" and "text" in event.get("content", {}).get("parts", [{}])[0]:
            assistant_message = event["content"]["parts"][0]["text"]
        
        # Look for text_to_speech function response to extract audio file path
        if "functionResponse" in event.get("content", {}).get("parts", [{}])[0]:
            func_response = event["content"]["parts"][0]["functionResponse"]
            if func_response.get("name") == "text_to_speech":
                response_text = func_response.get("response", {}).get("result", {}).get("content", [{}])[0].get("text", "")
                # Extract file path using simple string parsing
                if "File saved as:" in response_text:
                    parts = response_text.split("File saved as:")[1].strip().split()
                    if parts:
                        audio_file_path = parts[0].strip(".")
    
    # Add assistant response to chat
    if assistant_message:
        st.session_state.messages.append({"role": "assistant", "content": assistant_message, "audio_path": audio_file_path})
    
    return True

st.title("🧦 Sock Shop Vendor Portal")

with st.sidebar:
    st.header("Session Info")

    if st.session_state.session_id:
        st.success(f"Active session: {st.session_state.session_id}")
        if st.button("➕ New Session"):
            create_adk_session()
    else:
        st.warning("No active session")
        if st.button("➕ Create Session"):
            create_adk_session()

st.subheader("Conversation")
st.markdown("Welcome! Chat with our agent to learn how to add your socks to our store.")

for message in st.session_state.messages:
    with st.chat_message(message["role"]):
        st.markdown(message["content"])

if st.session_state.session_id:  # Only show input if session exists
    user_input = st.chat_input("Type your message...")
    if user_input:
        send_message(user_input)
        st.rerun()  # Rerun to update the UI with new messages
else:
    st.info("👈 Create a session to start chatting")

