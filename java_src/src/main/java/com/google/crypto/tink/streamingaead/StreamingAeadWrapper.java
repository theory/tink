// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

package com.google.crypto.tink.streamingaead;

import com.google.crypto.tink.PrimitiveSet;
import com.google.crypto.tink.PrimitiveWrapper;
import com.google.crypto.tink.Registry;
import com.google.crypto.tink.StreamingAead;
import java.security.GeneralSecurityException;
import java.util.ArrayList;
import java.util.List;

/**
 * StreamingAeadWrapper is the implementation of PrimitiveWrapper for the StreamingAead primitive.
 *
 * <p>The returned primitive works with a keyset (rather than a single key). To encrypt a plaintext,
 * it uses the primary key in the keyset. To decrypt, the primitive tries the enabled keys from the
 * keyset to select the right key for decryption. All keys in a keyset of StreamingAead have type
 * {@link com.google.crypto.tink.proto.OutputPrefixType#RAW}.
 */
public class StreamingAeadWrapper implements PrimitiveWrapper<StreamingAead, StreamingAead> {

  private static final StreamingAeadWrapper WRAPPER = new StreamingAeadWrapper();

  StreamingAeadWrapper() {}

  /**
   * @return a StreamingAead primitive from a {@code keysetHandle}.
   * @throws GeneralSecurityException
   */
  @Override
  public StreamingAead wrap(final PrimitiveSet<StreamingAead> primitives)
      throws GeneralSecurityException {
    List<StreamingAead> allStreamingAeads = new ArrayList<>();
    for (List<PrimitiveSet.Entry<StreamingAead>> entryList : primitives.getAll()) {
      // For legacy reasons (Tink always encrypted with non-RAW keys) we use all
      // primitives, even those which have output_prefix_type != RAW.
      for (PrimitiveSet.Entry<StreamingAead> entry : entryList) {
        allStreamingAeads.add(entry.getPrimitive());
      }
    }
    PrimitiveSet.Entry<StreamingAead> primary = primitives.getPrimary();
    if (primary == null || primary.getPrimitive() == null) {
      throw new GeneralSecurityException("No primary set");
    }
    return new StreamingAeadHelper(allStreamingAeads, primary.getPrimitive());
  }

  @Override
  public Class<StreamingAead> getPrimitiveClass() {
    return StreamingAead.class;
  }

  @Override
  public Class<StreamingAead> getInputPrimitiveClass() {
    return StreamingAead.class;
  }

  public static void register() throws GeneralSecurityException {
    Registry.registerPrimitiveWrapper(WRAPPER);
  }
}
