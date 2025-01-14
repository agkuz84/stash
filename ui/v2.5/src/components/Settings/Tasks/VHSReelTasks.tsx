import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { mutateGenerateVHSClips } from "src/core/StashService";

export const VHSReelTasks: React.FC = () => {
  const [isGenerating, setIsGenerating] = React.useState(false);
  const Toast = useToast();

  const generateClips = async () => {
    try {
      setIsGenerating(true);
      await mutateGenerateVHSClips();
      Toast.success("Started generating VHS clips");
    } catch (e) {
      Toast.error("Error starting VHS clip generation");
      console.error(e);
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div className="setting-section">
      <h1>VHS Reel Preparation</h1>
      <div className="setting">
        <div>
          <h3>Generate VHS Clips</h3>
          <div className="sub-heading">
            <FormattedMessage id="actions.generate_vhs_clips" defaultMessage="Creates 5-second clips for each marker in your library." />
            <FormattedMessage id="actions.generate_vhs_clips_required" defaultMessage="Required for VHS Reel playback." />
          </div>
        </div>
        <div>
          <Button
            variant="secondary"
            onClick={() => generateClips()}
            disabled={isGenerating}
          >
            {isGenerating ? (
              <LoadingIndicator message="Generating..." inline small />
            ) : (
              <>
                <Icon icon="tasks" />
                <span><FormattedMessage id="actions.generate" defaultMessage="Generate" /></span>
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
};